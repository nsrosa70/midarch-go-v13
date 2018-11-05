package components

import (
	"framework/messages"
)

type CalculatorClientProxy struct {
	Host  string
	Port  int
	Proxy string
}

var I_PreInvRCalculatorClientProxy = make(chan messages.SAMessage)
var I_PosTerRCalculatorClientProxy = make(chan messages.SAMessage)

func (c CalculatorClientProxy) Add(p1, p2 int) int {
	args := []int{p1, p2}
	inv := messages.Invocation{Host: c.Host, Port: c.Port, Op: "add", Args: args}
	reqMsg := messages.SAMessage{inv}

	I_PreInvRCalculatorClientProxy <- reqMsg
	repMsg := <-I_PosTerRCalculatorClientProxy
	payload := repMsg.Payload.(map[string]interface{})
	reply := int(payload["Reply"].(float64))
	return reply
}

func (CalculatorClientProxy) I_PreInvR(msg *messages.SAMessage,r *bool) {
	*msg = <-I_PreInvRCalculatorClientProxy
}

func (CalculatorClientProxy) I_PosTerR(msg *messages.SAMessage,r *bool) {
	I_PosTerRCalculatorClientProxy <- *msg
}