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
		go func() {
			defer done.Done()
			helloClient(requests)
		}()
	}
	done.Wait()
}

func helloClient(requests chan helloRequest) {
	request := newHelloRequest("Client")
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
