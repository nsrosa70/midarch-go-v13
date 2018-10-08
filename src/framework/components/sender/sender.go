package sender

import "framework/message"

type Sender struct{}


func (Sender) I_PreInvR(msg *message.Message) {
	*msg = message.Message{Payload:"testV2"}
}

func (e Sender) Loop(i_PreInvR, invR chan message.Message) {
	var msgIPreInvR message.Message
	for {
		select {
		case msgIPreInvR = <- i_PreInvR:
			e.I_PreInvR(&msgIPreInvR)
		case invR <- msgIPreInvR:
		}
	}
}
