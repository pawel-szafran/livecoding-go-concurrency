package main

import (
	"fmt"
	"sync"
)

func main() {
	var done sync.WaitGroup
	done.Add(1)
	go func() {
		defer done.Done()
		hello()
	}()
	done.Wait()
}

func hello() {
	fmt.Println("Hello Go :)")
}
