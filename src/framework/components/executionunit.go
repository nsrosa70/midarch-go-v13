package components

import (
	"framework/messages"
)

type ExecutionUnit struct {
}

var unitMsg messages.SAMessage

// Unit initialization
func (unit ExecutionUnit) I_InitialiseUnit(msg *messages.SAMessage, info *interface{}, r *bool) {
}

// Unit execution
func (unit ExecutionUnit) I_Execute(msg *messages.SAMessage, info *interface{}, r *bool) {
}

// Unit adaptation
func (unit ExecutionUnit) I_AdaptUnit(msg *messages.SAMessage, info *interface{}, r *bool) {
}
