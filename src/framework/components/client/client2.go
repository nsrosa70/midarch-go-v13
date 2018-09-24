package client

import (
	"framework/message"
	"fmt"
	"reflect"
)

type Client2 struct {}

func (Client2) I_PreInvR(msg *message.Message) {
	//time.Sleep(500 * time.Millisecond)
	*msg = message.Message{"client2"}
}

func (Client2) I_PosTerR(msg *message.Message) {
	fmt.Println(reflect.ValueOf(msg.Payload).String())
}