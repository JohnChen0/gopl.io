// Exercise 7.2 from page 174.

// Write a function CountingWriter that, given an io.Writer, returns a new
// Writer that wraps the original, and a pointer to an int64 variable that at
// any moment contains the number of bytes written to the new Writer.

package main

import (
	"fmt"
	"io"
	"os"
)

type writerWithCount struct {
	writer io.Writer
	count int64
}

func (w *writerWithCount) Write(p []byte) (int, error) {
	w.count += int64(len(p))
	return w.writer.Write(p)
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	result := &writerWithCount{w, 0}
	return result, &result.count
}

func main() {
	w, c := CountingWriter(os.Stdout)
	w.Write([]byte("hello"))
	fmt.Println(*c) // "5", = len("hello")

	*c = 0 // reset the counter
	var name = "Dolly"
	fmt.Fprintf(w, "hello, %s", name)
	fmt.Println(*c) // "12", = len("hello, Dolly")
}
