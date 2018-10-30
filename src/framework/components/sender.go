package components

import (
	"framework/messages"
	"strconv"
)

type Sender struct{}

var idx = 0

func (Sender) I_PreInvR1(msg *messages.SAMessage) {
	*msg = messages.SAMessage{Payload:"Message 01 ["+strconv.Itoa(idx)+"]"}
	idx ++
}

func (Sender) I_PreInvR2(msg *messages.SAMessage) {
	*msg = messages.SAMessage{Payload:"Message 02 ["+strconv.Itoa(idx)+"]"}
	idx ++
}
