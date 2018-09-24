package fibonacciclientproxy

import (
	"framework/message"
)

type FibonacciClientProxy struct {
	Host  string
	Port  int
	Proxy string
}

var chIn = make(chan message.Message)
var chOut = make(chan message.Message)

func (c FibonacciClientProxy) Fibo(p1 int) int {
	args := []int{p1}
	inv := message.Invocation{Host: c.Host, Port: c.Port, Op: "fibo", Args: args}
	reqMsg := message.Message{inv}

	chIn <- reqMsg
	repMsg := <-chOut
	payload := repMsg.Payload.(map[string]interface{})
	reply := int(payload["Reply"].(float64))
	return reply
}

func (FibonacciClientProxy) I_PreInvR(msg *message.Message) {
	*msg = <-chIn
}

func (FibonacciClientProxy) I_PosTerR(msg *message.Message) {
	chOut <- *msg
}