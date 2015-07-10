package main

import "fmt"

func main() {
	hellos := make(chan string)
	go func() {
		hellos <- hello()
	}()
	fmt.Println(<-hellos)
}

func hello() string {
	return "Hello Go :)"
}
