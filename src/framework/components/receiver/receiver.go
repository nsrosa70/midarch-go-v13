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

func (Receiver) Loop(channels map[string]chan message.Message) {
	var msgIPosInvP message.Message
	for {
		select {
		case <-channels["InvP"]:
		case msgIPosInvP = <-channels["I_PosInvP_receiver"]:
			Receiver{}.I_PosInvP(&msgIPosInvP)
		}
	}
}