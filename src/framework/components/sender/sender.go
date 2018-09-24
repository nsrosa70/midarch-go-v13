package sender

import (
	"framework/message"
	"time"
)

type Sender struct {}

func (Sender) I_PreInvR(msg *message.Message) {
	time.Sleep(1 * time.Second)
	*msg = message.Message{"test"}
}
