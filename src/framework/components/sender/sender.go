package sender

import (
	"strconv"
	"fmt"
)

type Sender struct{}

func (Sender) T() {
	fmt.Println("Here")
}

//func (Sender) Loop(invR, terR, invP, terP, i_PosInvP chan string) {
func (Sender) Loop(invR chan string) {
	msg := "testV1"
	i := 0
	//msgRecv := ""
	for {
		select {
		case invR <- msg + strconv.Itoa(i):
			//fmt.Println("Sender:: invR")
		//case terP <- msgRecv:
			//fmt.Println("Sender:: terP")
		//case msgRecv = <-terR:
			//fmt.Println("Sender:: terR")
		//case <-invP:
			//fmt.Println("Sender:: invP")
		}
	}
}
