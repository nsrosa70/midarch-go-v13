package main

import (
	"framework/message"
	"apps/fibonacci/impl"
	"fmt"
)

type FibonacciInvoker struct{}

var msg message.Message

func GetTypeElem() interface{}{
	return FibonacciInvoker{}
}

func GetBehaviourExp() string {
	//return library.BehaviourLibrary[calculatorinvoker.CalculatorInvoker{}]
	return "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"
}

func (FibonacciInvoker) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.MIOP).RequestBody.Op
	args := msg.Payload.(message.MIOP).RequestBody.Args

	switch op {
	case "fibo":
		// prepare invocation
		argsX := args.([]interface{})
		p1 := int(argsX[0].(float64))

		r := impl.Fibo(p1)
		fmt.Println("[PLUGIN 01]")
		fmt.Println(r)

		// send reply
		header := message.ReplyHeader{1} // 1 - Success
		body := message.ReplyBody{Reply: r}
		miop := message.MIOP{ReplyHeader: header, ReplyBody: body}
		*msg = message.Message{miop}
	}
}