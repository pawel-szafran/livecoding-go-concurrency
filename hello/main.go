package main

import (
	"fmt"
	"time"
)

func main() {
	go hello()
	time.Sleep(time.Second)
}

func hello() {
	fmt.Println("Hello Go :)")
}
