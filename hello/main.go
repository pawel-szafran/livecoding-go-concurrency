package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	requests := make(chan helloRequest)
	go helloBroker(requests)
	lotsOfHelloClients(requests)
	close(requests)
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

func lotsOfHelloClients(requests chan helloRequest) {
	var done sync.WaitGroup
	for i := 0; i < 10000; i++ {
		done.Add(1)
		go func(id int) {
			defer done.Done()
			helloClient(id, requests)
		}(i)
	}
	done.Wait()
}

func helloClient(id int, requests chan helloRequest) {
	request := newHelloRequest(fmt.Sprint("Client-", id))
	requests <- request
	fmt.Println(<-request.response)
}

func helloBroker(requests chan helloRequest) {
	for i := 0; i < 1000; i++ {
		go helloWorker(requests)
	}
}

func helloWorker(requests chan helloRequest) {
	for request := range requests {
		request.response <- hello(request.name)
	}
}

func hello(name string) string {
	time.Sleep(50 * time.Millisecond)
	return fmt.Sprint("Hello ", name, " :)")
}
