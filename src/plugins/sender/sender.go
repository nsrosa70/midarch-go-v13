package main

import (
	"framework/message"
	"time"
	"framework/library"
)

type Sender struct{}

var msg message.Message

func GetTypeElem() interface{}{
	return Sender{}
}

func GetBehaviourExp() string {
	return library.SENDER_BEHAVIOUR
}

func (Sender) I_PreInvR(msg *message.Message) {
	time.Sleep(1 * time.Second)
	*msg = message.Message{"WELLISON"}
}
