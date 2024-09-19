// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// +build ignore

// The thumbnail command produces thumbnails of JPEG files
// whose names are provided on each line of the standard input.
//
// The "+build ignore" tag (see p.295) excludes this file from the
// thumbnail package, but it can be compiled as a command and run like
// this:
//
// Run with:
//   $ go run $GOPATH/src/gopl.io/ch8/thumbnail/main.go
//   foo.jpeg
//   ^D
//
// This file has been modified from the original to provide the capability to
// call the example functions makeThumbnails* from the book. For example, to
// call makeThumbnails5, use "go run main.go 5". I have also added my own
// version makeThumbnails7, which can be invoked with "go run main.go 7".
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"gopl.io/ch8/thumbnail"
)

func usage() {
	fmt.Printf("usage: %s [0-7]\n", os.Args[0])
	os.Exit(1)
}

func parseCommandLine() int {
	switch len(os.Args) {
	case 1:
		return 0
	case 2:
		arg := os.Args[1]
		if len(arg) != 1 || arg[0] < '0' || arg[0] > '7' {
			usage()
			panic("should not reach")
		}
		return int(arg[0]) - int('0')
	default:
		usage()
		panic("should not reach")
	}
}

func main() {
	funcIndex := parseCommandLine()
	var filenamesSlice []string
	var filenamesChan chan string
	var sizeChan chan int64
	if funcIndex >= 6 {
		filenamesChan = make(chan string)
		sizeChan = make(chan int64)
		go func() {
			if funcIndex == 6 {
				sizeChan <- thumbnail.MakeThumbnails6(filenamesChan)
			} else {
				sizeChan <- thumbnail.MakeThumbnails7(filenamesChan)
			}
		}()
	}
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		if funcIndex >= 6 {
			filenamesChan <- input.Text()
			continue
		}
		if funcIndex > 0 {
			filenamesSlice = append(filenamesSlice, input.Text())
			continue
		}
		thumb, err := thumbnail.ImageFile(input.Text())
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println(thumb)
	}
	if err := input.Err(); err != nil {
		log.Fatal(err)
	}
	switch funcIndex {
	case 1:
		thumbnail.MakeThumbnails(filenamesSlice)
	case 2:
		thumbnail.MakeThumbnails2(filenamesSlice)
	case 3:
		thumbnail.MakeThumbnails3(filenamesSlice)
	case 4:
		err := thumbnail.MakeThumbnails4(filenamesSlice)
		if err != nil {
			log.Fatal(err)
		}

	case 5:
		thumbfiles, err := thumbnail.MakeThumbnails5(filenamesSlice)
		if err != nil {
			log.Fatal(err)
		} else {
			for _, f := range thumbfiles {
				fmt.Println(f)
			}
		}
	case 6, 7:
		close(filenamesChan)
		fmt.Println("Total size:", <-sizeChan)
	}
}
