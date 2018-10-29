package components

import (
	"fmt"
	"framework/messages"
)

type Receiver struct{}

func (Receiver) I_PosInvP(m *messages.SAMessage) {
	fmt.Print("Receiver::::::::::::::::::::::: ")
	fmt.Println(m.Payload)
}