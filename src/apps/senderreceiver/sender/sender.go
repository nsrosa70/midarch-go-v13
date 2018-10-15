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

func (e Sender) Loop(channels map[string] chan message.Message) {
	var msgIPreInvR message.Message

	for {
		select {
		case msgIPreInvR = <- channels["I_PreInvR_sender"]:
			e.I_PreInvR(&msgIPreInvR)
		case channels["InvR"] <- msgIPreInvR:
		}
	}
}

