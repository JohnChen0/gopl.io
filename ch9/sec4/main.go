package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	var x, y int
	go func() {
		x = 1			// A1
		fmt.Print("y:", y, " ")	// A2
		wg.Done()
	}()
	go func() {
		y = 1			// B1
		fmt.Print("x:", x, " ")	// B2
		wg.Done()
	}()
	wg.Wait()
	fmt.Println()
}
