package receiver

import (
	"fmt"
	"framework/message"
	"reflect"
)

type Receiver struct{}

func (Receiver) I_PosInvP(msg *message.Message) {
	fmt.Println("Received from Sender: " + reflect.ValueOf(msg.Payload).String())
}
