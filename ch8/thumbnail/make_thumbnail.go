// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// This file is just a place to put example code from the book.
// It does not actually run any code in gopl.io/ch8/thumbnail.

// The following changes have been made to this file from the original:
// * Renamed and moved to thumbnail package, so its functions can be called.
// * The makeThumbnails* functions are now exported.
// * Added MakeThumbnails7 function, which is similar to MakeThumbnails6 from
//   the book, but without using WaitGroup.

package thumbnail

import (
	"log"
	"os"
	"sync"
)

//!+1
// makeThumbnails makes thumbnails of the specified files.
func MakeThumbnails(filenames []string) {
	for _, f := range filenames {
		if _, err := ImageFile(f); err != nil {
			log.Println(err)
		}
	}
}

//!-1

//!+2
// NOTE: incorrect!
func MakeThumbnails2(filenames []string) {
	for _, f := range filenames {
		go ImageFile(f) // NOTE: ignoring errors
	}
}

//!-2

//!+3
// makeThumbnails3 makes thumbnails of the specified files in parallel.
func MakeThumbnails3(filenames []string) {
	ch := make(chan struct{})
	for _, f := range filenames {
		go func(f string) {
			ImageFile(f) // NOTE: ignoring errors
			ch <- struct{}{}
		}(f)
	}

	// Wait for goroutines to complete.
	for range filenames {
		<-ch
	}
}

//!-3

//!+4
// makeThumbnails4 makes thumbnails for the specified files in parallel.
// It returns an error if any step failed.
func MakeThumbnails4(filenames []string) error {
	errors := make(chan error)

	for _, f := range filenames {
		go func(f string) {
			_, err := ImageFile(f)
			errors <- err
		}(f)
	}

	for range filenames {
		if err := <-errors; err != nil {
			return err // NOTE: incorrect: goroutine leak!
		}
	}

	return nil
}

//!-4

//!+5
// makeThumbnails5 makes thumbnails for the specified files in parallel.
// It returns the generated file names in an arbitrary order,
// or an error if any step failed.
func MakeThumbnails5(filenames []string) (thumbfiles []string, err error) {
	type item struct {
		thumbfile string
		err       error
	}

	ch := make(chan item, len(filenames))
	for _, f := range filenames {
		go func(f string) {
			var it item
			it.thumbfile, it.err = ImageFile(f)
			ch <- it
		}(f)
	}

	for range filenames {
		it := <-ch
		if it.err != nil {
			return nil, it.err
		}
		thumbfiles = append(thumbfiles, it.thumbfile)
	}

	return thumbfiles, nil
}

//!-5

//!+6
// makeThumbnails6 makes thumbnails for each file received from the channel.
// It returns the number of bytes occupied by the files it creates.
func MakeThumbnails6(filenames <-chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup // number of working goroutines
	for f := range filenames {
		wg.Add(1)
		// worker
		go func(f string) {
			defer wg.Done()
			thumb, err := ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}
			info, _ := os.Stat(thumb) // OK to ignore error
			sizes <- info.Size()
		}(f)
	}

	// closer
	go func() {
		wg.Wait()
		close(sizes)
	}()

	var total int64
	for size := range sizes {
		total += size
	}
	return total
}

//!-6

// MakeThumbnails7 is similar to MakeThumbnails6, but without using WaitGroup.
func MakeThumbnails7(filenames <-chan string) int64 {
	sizes := make(chan int64)
	remaining := 0
	for f := range filenames {
		remaining++
		// worker
		go func(f string) {
			thumb, err := ImageFile(f)
			if err != nil {
				log.Println(err)
				sizes <- 0 // Making sure every worker sends a size
				return
			}
			info, _ := os.Stat(thumb) // OK to ignore error
			sizes <- info.Size()
		}(f)
	}

	var total int64
	for remaining > 0 {
		total += <-sizes
		remaining--
	}
	return total
}
