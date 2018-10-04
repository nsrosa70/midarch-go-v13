package receiver

import (
	"fmt"
)

type Receiver struct{}

func (Receiver) I_PosInvP(m *string) {
	fmt.Println("Receiver::::::::::::::::::::::: " + *m)
}

//func (Receiver) Loop(invR, terR, invP, terP, i_PosInvP chan string) {
func (Receiver) Loop(invP, i_PosInvP chan string) {
	//msg := "testV1"
	//i := 0
	//msgRcv := ""
	for {
		select {
		//case invR <- msg + strconv.Itoa(i):
			//fmt.Println("Receiver:: invR")
		//	i++
		//case terP <- msgRcv:
			//fmt.Println("Receiver:: terP")
		//case msgRcv = <-terR:
			//fmt.Println("Receiver:: terR")
		case <-invP:
			//fmt.Println("receiver:: InvP")
		case msgRcv := <-i_PosInvP:
			//fmt.Println("Receiver:: i_PosInvP")
			Receiver{}.I_PosInvP(&msgRcv)
		}
	}
}