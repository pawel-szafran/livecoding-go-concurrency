package main

import (
	"fmt"
	"time"
)

func main() {
	requests := make(chan chan string)
	go helloBroker(requests)
	helloClient(requests)
}

func helloClient(requests chan chan string) {
	response := make(chan string)
	requests <- response
	fmt.Println(<-response)
}

func helloBroker(requests chan chan string) {
	go helloWorker(requests)
}

func helloWorker(requests chan chan string) {
	response := <-requests
	response <- hello()
}

func hello() string {
	time.Sleep(50 * time.Millisecond)
	return "Hello Go :)"
}
