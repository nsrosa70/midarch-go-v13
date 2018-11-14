package components

import (
	"framework/messages"
	"fmt"
)

type MAPEKExecutor struct {}

func (MAPEKExecutor) I_PosInvP(msg *messages.SAMessage, r *bool) {
	fmt.Println(msg.Payload)
}