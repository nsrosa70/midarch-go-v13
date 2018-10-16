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