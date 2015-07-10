package main

import "fmt"

func main() {
	done := make(chan bool)
	go func() {
		defer func() {
			done <- true
		}()
		hello()
	}()
	<-done
}

func hello() {
	fmt.Println("Hello Go :)")
}
