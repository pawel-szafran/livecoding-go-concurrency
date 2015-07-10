package main

import "fmt"

func main() {
	hellos := make(chan string)
	go func() {
		hellos <- hello()
	}()
	hello := <-hellos
	fmt.Println(hello)
}

func hello() string {
	return "Hello Go :)"
}
