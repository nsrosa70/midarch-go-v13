package sender

import (
	"framework/message"
	"strconv"
)

type Sender struct{}

var idx = 0

func (Sender) I_PreInvR(msg *message.Message) {
	*msg = message.Message{Payload:"Message sent ["+strconv.Itoa(idx)+"]"}
	idx ++
}

