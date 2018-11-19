package components

import (
	"framework/messages"
	"fmt"
)

type Receiver struct{}

func (Receiver) I_PosInvP(msg *messages.SAMessage, r *bool) {
	fmt.Println("Receiver:: " + msg.Payload.(string))
	return
}
