package components

import (
	"fmt"
	"gmidarch/development/framework/messages"
)

type Receiver struct{}

func (Receiver) I_PosInvP(msg *messages.SAMessage, info interface{}, r *bool) {
	//fmt.Println("Receiver:: " + msg.Payload.(string))
	fmt.Printf("Receiver:: %v\n",msg.Payload)
	return
}
