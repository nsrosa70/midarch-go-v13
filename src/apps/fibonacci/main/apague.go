package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	for {
		select {
		case c <- timeout(c):
			fmt.Println("timeout")
			break
		default:
			fmt.Println("loop")
		}
	}
}

func timeout(c chan int) int {
	time.Sleep(1 * time.Millisecond)

	<- c
	return 1
}
