package client

import (
	"framework/message"
	"fmt"
	"reflect"
)

type Client1 struct {}

func (Client1) I_PreInvR(msg *message.Message) {
	//time.Sleep(500 * time.Millisecond)
	*msg = message.Message{"client1"}
}

func (Client1) I_PosTerR(msg *message.Message) {
	fmt.Println(reflect.ValueOf(msg.Payload).String())
}