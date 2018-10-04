package client

import "fmt"

type Client struct{}

func (Client) Loop(invR,terR chan string) {
	msg := "testV1"
	for {
		select {
		case invR <- msg :
			//fmt.Println("Sender:: invR")
			//case terP <- msgRecv:
			//fmt.Println("Sender:: terP")
			case msgRecv := <-terR:
			fmt.Println("Sender:: terR "+msgRecv)
			//case <-invP:
			//fmt.Println("Sender:: invP")
		}
	}
}

