package server

import (
	"framework/message"
)

type Server struct {}

func (Server) I_PosInvP(msg *message.Message) {
	*msg = message.Message{Payload:msg.Payload}
}
