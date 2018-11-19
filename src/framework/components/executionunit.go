package components

import (
	"framework/messages"
	"fmt"
)

type ExecutionUnit struct{}

func (ExecutionUnit) I_Execute(msg *messages.SAMessage, r *bool) {
	fmt.Printf("Unit:: %v \n",msg.Payload)
}
