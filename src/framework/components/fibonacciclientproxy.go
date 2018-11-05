package components

import (
	"framework/messages"
	"shared/parameters"
)

type FibonacciClientProxy struct {
	Host string
	Port int
}

var i_PreInvRFibonacciClientProxy = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
var i_PosTerTFibonacciClientProxy = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)

func (c FibonacciClientProxy) Fibo(_p1 int) int {
	c.Port = parameters.FIBONACCI_PORT // TODO
	_args := []int{_p1}
	reqMsg := messages.SAMessage{messages.Invocation{Host: c.Host, Port: c.Port, Op: "Fibo", Args: _args}}

	i_PreInvRFibonacciClientProxy <- reqMsg
	repMsg := <-i_PosTerTFibonacciClientProxy

	payload := repMsg.Payload.(map[string]interface{})
	reply := int(payload["Reply"].(float64))

	return reply
}

func (FibonacciClientProxy) I_PreInvR(msg *messages.SAMessage, r *bool) {
	*msg = <-i_PreInvRFibonacciClientProxy
}

func (FibonacciClientProxy) I_PosTerR(msg *messages.SAMessage, r *bool) {
	i_PosTerTFibonacciClientProxy <- *msg
}
