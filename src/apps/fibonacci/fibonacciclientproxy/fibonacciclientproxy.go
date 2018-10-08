package fibonacciclientproxy

import (
	"framework/message"
	"shared/parameters"
)

type FibonacciClientProxy struct {
	Host  string
	Port int
}

var i_PreInvR  = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
var i_PosTerT = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)

func (e FibonacciClientProxy) Loop(I_PreInvR, InvR, TerR, I_PosTerR chan message.Message) {
	var msgTerR, msgPreInvR message.Message
	for {
		select {
		case msgPreInvR = <-I_PreInvR:
			e.I_PreInvR(&msgPreInvR)
		case InvR <- msgPreInvR:
		case msgTerR = <-TerR:
		case <-I_PosTerR:
			e.I_PosTerR(&msgTerR)
		}
	}
}

func (c FibonacciClientProxy) Fibo(_p1 int) int {
	c.Port = parameters.FIBONACCI_PORT // TODO
	_args := []int{_p1}
	reqMsg := message.Message{message.Invocation{Host: c.Host, Port: c.Port, Op: "Fibo", Args: _args}}

	i_PreInvR  <- reqMsg
	repMsg := <-i_PosTerT

	payload := repMsg.Payload.(map[string]interface{})
	reply := int(payload["Reply"].(float64))

	return reply
}

func (FibonacciClientProxy) I_PreInvR(msg *message.Message) {
	*msg = <-i_PreInvR
}

func (FibonacciClientProxy) I_PosTerR(msg *message.Message) {
	i_PosTerT <- *msg
}