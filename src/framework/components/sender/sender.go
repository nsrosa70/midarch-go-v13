package sender

import "framework/message"

type Sender struct{}


func (Sender) I_PreInvR(msg *message.Message) {
	*msg = message.Message{Payload:"testV2"}
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
