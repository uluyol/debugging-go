package main

import (
	"time"
)

func computeLoop(c <-chan int) {
	select {
	case <-c:
		return
	case <-time.After(2 * time.Second):
		panic("should have sent by now")
	}
}

func main() {
	c1 := make(chan int)
	c2 := make(chan int)
	go computeLoop(c1)
	c2 <- 1
}
