package components

import (
	"framework/messages"
)

type MAPEKExecutor struct {}

func (MAPEKExecutor) I_Execute(msg *messages.SAMessage, r *bool) {
	//fmt.Printf("Executor:: %i \n",msg.Payload)
}