package main

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan chan string)
	go helloWorker(requests)
	response := make(chan string)
	requests <- response
	fmt.Println(<-response)
}

func helloWorker(requests chan chan string) {
	response := <-requests
	response <- hello()
}

func hello() string {
	time.Sleep(50 * time.Millisecond)
	return "Hello Go :)"
}
