package components

import (
	"framework/messages"
)

type ExecutionEnvironment struct{}

func (ExecutionEnvironment) I_DefineUnit(msg *messages.SAMessage, r *bool) {
	*msg = messages.SAMessage{Payload:"Unit 1"}
}