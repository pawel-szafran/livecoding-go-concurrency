package main

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan chan string)
	go func() {
		response := <-requests
		response <- hello()
	}()
	response := make(chan string)
	requests <- response
	fmt.Println(<-response)
}

func hello() string {
	time.Sleep(50 * time.Millisecond)
	return "Hello Go :)"
}
