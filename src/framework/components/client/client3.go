package client

import (
	"framework/message"
	"fmt"
	"reflect"
)

type Client3 struct {}

func (Client3) I_PreInvR(msg *message.Message) {
	//time.Sleep(500 * time.Millisecond)
	*msg = message.Message{"client3"}
}

func (Client3) I_PosTerR(msg *message.Message) {
	fmt.Println(reflect.ValueOf(msg.Payload).String())
}