package main

import (
	"framework/message"
	"apps/fibonacci/fibonacci"
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

func (n FibonacciInvoker) Loop(channels map[string]chan message.Message) {
	var msgPosInvP message.Message
	for {
		select {
		case <-channels["InvP"]:
		case msgPosInvP = <-channels["I_PosInvP_fibonacciinvoker"]:
			n.I_PosInvP(&msgPosInvP)
		case channels["TerP"] <- msgPosInvP:
			return
		}
	}
}

func (FibonacciInvoker) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.MIOP).Body.RequestHeader.Operation

	switch op {
	case "Fibo":
		// process request
		_args := msg.Payload.(message.MIOP).Body.RequestBody.Args
		_argsX := _args.([]interface{})
		_p1 := int(_argsX[0].(float64))
		_r := fibonacci.Fibonacci{}.Fibo(_p1) // dispatch

		//fmt.Println("Plugin 01")

		// send reply
		_replyHeader := message.ReplyHeader{Status: 1} // 1 - Success
		_replyBody := message.ReplyBody{Reply: _r}
		_miopHeader := message.MIOPHeader{Magic: "MIOP"}
		_miopBody := message.MIOPBody{ReplyHeader: _replyHeader, ReplyBody: _replyBody}
		_miop := message.MIOP{Header: _miopHeader, Body: _miopBody}
		*msg = message.Message{_miop}
	default:
		fmt.Println("FIBONACCIINVOKER:: Operation " + op + " not supported")
	}
}
