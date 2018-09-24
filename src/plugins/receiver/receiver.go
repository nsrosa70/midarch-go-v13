package main

import (
	"framework/message"
	"fmt"
	"framework/library"
)

type Receiver struct{}

var msg message.Message

func GetTypeElem() interface{}{
	return Receiver{}
}

func GetBehaviourExp() string {
	return library.RECEIVER_BEHAVIOUR
}

func (Receiver) I_PosInvP(msg *message.Message) {
	fmt.Println("[plugin Receiver] Received from Sender: " + msg.Payload)
}
