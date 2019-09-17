package main

import (
	"fmt"
)

func main() {
	chn1 := make(chan string)
	//chn2 := make(chan string)

	go func() {
		for {
			chn1 <- "from 1"
			//time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		for {
			chn1 <- "from 2"
			//time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		for {
			select {
			case chn1 <- "1":
				//fmt.Println(msg1)
			case chn1 <- "1":
			//msg2 := <-chn1:
			//	fmt.Println(msg2)
			default:
				fmt.Printf("Não está havendo sincronização\n")
			}
		}
	}()
	fmt.Scanln()
}
