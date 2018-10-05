package sender

import (
	"fmt"
)

type Sender struct{}

func (Sender) T() {
	fmt.Println("Here")
}

func (Sender) Loop(invR chan string) {
	msg := "testV1"
	for {
		select {
		case invR <- msg:
		}
	}
}
