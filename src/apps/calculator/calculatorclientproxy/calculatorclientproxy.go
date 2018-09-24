package calculatorclientproxy

import (
	"framework/message"
)

type CalculatorClientProxy struct {
	Host  string
	Port  int
	Proxy string
}

var chIn = make(chan message.Message)
var chOut = make(chan message.Message)

func (c CalculatorClientProxy) Add(p1, p2 int) int {
	args := []int{p1, p2}
	inv := message.Invocation{Host: c.Host, Port: c.Port, Op: "add", Args: args}
	reqMsg := message.Message{inv}

	chIn <- reqMsg
	repMsg := <-chOut
	payload := repMsg.Payload.(map[string]interface{})
	reply := int(payload["Reply"].(float64))
	return reply
}

func (CalculatorClientProxy) I_PreInvR(msg *message.Message) {
	*msg = <-chIn
}

func (CalculatorClientProxy) I_PosTerR(msg *message.Message) {
	chOut <- *msg
}

