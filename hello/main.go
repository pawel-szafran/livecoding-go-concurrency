package main

import (
	"fmt"
	"time"
)

func main() {
	hellos := make(chan string)
	go func() {
		hellos <- hello()
	}()
	fmt.Println(<-hellos)
}

func hello() string {
	time.Sleep(50 * time.Millisecond)
	return "Hello Go :)"
}
