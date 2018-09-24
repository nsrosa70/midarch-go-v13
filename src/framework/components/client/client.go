package client

import (
	"framework/message"
	"fmt"
	"strconv"
)

type Client struct {}

var count int

func (Client) I_PreInvR(msg *message.Message) {
	//time.Sleep(500 * time.Millisecond)

	header := message.RequestHeader{Magic:"MIOP"}
	count++
	body := message.RequestBody{Op:strconv.Itoa(count)}
	miop := message.MIOP{RequestHeader:header,RequestBody:body}
	toCRH := message.ToCRH{Host:"localhost",Port:7070,MIOP:miop}

	//*msg := message.Message{"teste"}
	*msg = message.Message{toCRH}
}

func (Client) I_PosTerR(msg *message.Message) {
	payload := msg.Payload
	fmt.Println(payload)
	//fmt.Println(reflect.ValueOf(msg.Payload))
}