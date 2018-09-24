package fibonacciinvoker

import (
	"framework/message"
	"apps/fibonacci/impl"
)

type FibonacciInvoker struct{}

func (FibonacciInvoker) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.MIOP).RequestBody.Op
	args := msg.Payload.(message.MIOP).RequestBody.Args

	switch op {
	case "fibo":
		// prepare invocation
		argsX := args.([]interface{})
		p1 := int(argsX[0].(float64))

		//dispatch invocation
		r := impl.Fibo(p1)

		// send reply
		header := message.ReplyHeader{1} // 1 - Success
		body := message.ReplyBody{Reply: r}
		miop := message.MIOP{ReplyHeader: header, ReplyBody: body}
		*msg = message.Message{miop}
	}
}
