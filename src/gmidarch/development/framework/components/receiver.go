package components

import (
	"gmidarch/development/framework/messages"
	"fmt"
)

type Receiver struct{}

func (Receiver) I_PosInvP(msg *messages.SAMessage, info interface{}, r *bool) {
	//fmt.Println("Receiver:: " + msg.Payload.(string))
	fmt.Printf("Receiver:: %v\n",msg.Payload)
	return
}
