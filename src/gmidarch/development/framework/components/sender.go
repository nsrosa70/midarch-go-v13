package components

import (
	"strconv"
	"time"
	"gmidarch/development/framework/messages"
)

type Sender struct{}

var idx1 = 0
var idx2 = 0

func (Sender) I_PreInvR1(msg *messages.SAMessage, info interface{}, r *bool) {
	*msg = messages.SAMessage{Payload:"Message 01 ["+strconv.Itoa(idx1)+"]"}
	idx1++

	time.Sleep(1*time.Millisecond)
	*r = true
	return
}

func (Sender) I_PreInvR2(msg *messages.SAMessage, info interface{}, r *bool) {
	*msg = messages.SAMessage{Payload:"Message 02 ["+strconv.Itoa(idx2)+"]"}
	idx2++
	time.Sleep(1*time.Millisecond)
	*r = true
	return
}
