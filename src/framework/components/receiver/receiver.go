package receiver

import (
	"fmt"
	"framework/message"
)

type Receiver struct{}

func (Receiver) I_PosInvP(m *message.Message) {
	fmt.Print("Receiver::::::::::::::::::::::: ")
	fmt.Println(m.Payload)
}

func (Receiver) Loop(invP, i_PosInvP chan message.Message) {
	for {
		select {
		case <-invP:
		case msgRcv := <-i_PosInvP:
			Receiver{}.I_PosInvP(&msgRcv)
		}
	}
}