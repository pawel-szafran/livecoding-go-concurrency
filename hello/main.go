package main

import "fmt"

func main() {
	done := make(chan bool)
	go func() {
		hello()
		done <- true
	}()
	<-done
}

func hello() {
	fmt.Println("Hello Go :)")
}
