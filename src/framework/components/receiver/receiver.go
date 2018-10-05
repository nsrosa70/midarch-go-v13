package receiver

import (
	"fmt"
)

type Receiver struct{}

func (Receiver) I_PosInvP(m *string) {
	fmt.Println("Receiver::::::::::::::::::::::: " + *m)
}

func (Receiver) Loop(invP, i_PosInvP chan string) {
	for {
		select {
		case <-invP:
		case msgRcv := <-i_PosInvP:
			//fmt.Println("Receiver:: i_PosInvP")
			Receiver{}.I_PosInvP(&msgRcv)
		}
	}
}