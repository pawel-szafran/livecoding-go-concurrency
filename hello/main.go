package main

import (
	"fmt"
	"sync"
)

func main() {
	var done sync.WaitGroup
	done.Add(1)
	go func() {
		hello()
		done.Done()
	}()
	done.Wait()
}

func hello() {
	fmt.Println("Hello Go :)")
}
