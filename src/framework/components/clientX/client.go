package clientX

import "fmt"

type ClientX struct{}

func (ClientX) Loop(invR, terR chan string) {
	msgSent := "testV1"
	for {
		select {
		case invR <- msgSent:
		case msgReceived := <-terR:
			fmt.Println("Client:: terR :: " + msgReceived)
		}
	}
}
