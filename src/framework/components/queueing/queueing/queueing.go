package queueing

import (
	"framework/message"
	"fmt"
	"framework/components/queueing/queueimpl"
)

type QueueingService struct{}

var qs = queueimpl.QueueServiceImpl{}

func (QueueingService) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.Invocation).Op

	switch op {
	case "publish":
		args := msg.Payload.(message.Invocation).Args
		argsX := args.([]interface{})
		_p1 := argsX[0].(string)
		_p2 := argsX[1].(string)
		_r := qs.Publish(_p1, _p2)
		msgRep := message.Message{Payload: _r}
		*msg = msgRep
	case "consume":
		args := msg.Payload.(message.Invocation).Args
		argsX := args.([]interface{})
		_p1 := argsX[0].(string)
		_r := qs.Consume(_p1)
		msgRep := message.Message{Payload: _r}
		*msg = msgRep
	case "subscribe":
	default:
		fmt.Println("Queueing:: Operation " + op + " is not implemented by Queue Server")
	}
}

