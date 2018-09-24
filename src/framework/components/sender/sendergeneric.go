package sender

import (
	"framework/message"
	"time"
)

type SenderGeneric struct {}

func (SenderGeneric) I_PreInvR(msg *message.Message) {
	time.Sleep(1 * time.Second)
	receive()
	*msg = message.Message{"test"}
}

func receive(){

}