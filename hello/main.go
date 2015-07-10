package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	requests := make(chan helloRequest)
	go helloBroker(requests)
	helloClient(requests)
	elapsed := time.Since(start)
	fmt.Println("Elapsed:", elapsed)
}

type helloRequest struct {
	name     string
	response chan string
}

func newHelloRequest(name string) helloRequest {
	return helloRequest{name, make(chan string)}
}

func helloClient(requests chan helloRequest) {
	request := newHelloRequest("Client")
	requests <- request
	fmt.Println(<-request.response)
}

func helloBroker(requests chan helloRequest) {
	go helloWorker(requests)
}

func helloWorker(requests chan helloRequest) {
	request := <-requests
	request.response <- hello(request.name)
}

func hello(name string) string {
	time.Sleep(50 * time.Millisecond)
	return fmt.Sprint("Hello ", name, " :)")
}
