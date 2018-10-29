package components

import (
	"framework/messages"
	"strconv"
)

type Sender struct{}

var idx = 0

func (Sender) I_PreInvR(msg *messages.SAMessage) {
	*msg = messages.SAMessage{Payload:"Message sent ["+strconv.Itoa(idx)+"]"}
	idx ++
}

