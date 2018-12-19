package main

import (
	"framework/messages"
	"apps/fibonacci/fibonacci"
	"fmt"
)

type FibonacciInvoker struct{}

var msg messages.SAMessage

func GetTypeElem() interface{}{
	return FibonacciInvoker{}
}

func GetBehaviourExp() string {
	//return libraries.BehaviourLibrary[calculatorinvoker.CalculatorInvoker{}]
	return "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"
}

func (n FibonacciInvoker) Loop(channels map[string]chan messages.SAMessage) {
	var msgPosInvP messages.SAMessage
	for {
		select {
		case <-channels["InvP"]:
		case msgPosInvP = <-channels["I_PosInvP"]:
			n.I_PosInvP(&msgPosInvP)
		case channels["TerP"] <- msgPosInvP:
			return
		}
	}
}

func (FibonacciInvoker) I_PosInvP(msg *messages.SAMessage) {
	op := msg.Payload.(messages.MIOP).Body.RequestHeader.Operation

	switch op {
	case "Fibo":
		// process request
		_args := msg.Payload.(messages.MIOP).Body.RequestBody.Args
		_argsX := _args.([]interface{})
		_p1 := int(_argsX[0].(float64))
		_r := fibonacci.Fibonacci{}.Fibo(_p1) // dispatch

		fmt.Println("Plugin 02")

		// send reply
		_replyHeader := messages.ReplyHeader{Status: 1} // 1 - Success
		_replyBody := messages.ReplyBody{Reply: _r}
		_miopHeader := messages.MIOPHeader{Magic: "MIOP"}
		_miopBody := messages.MIOPBody{ReplyHeader: _replyHeader, ReplyBody: _replyBody}
		_miop := messages.MIOP{Header: _miopHeader, Body: _miopBody}
		*msg = messages.SAMessage{_miop}
	default:
		fmt.Println("FIBONACCIINVOKER:: Operation " + op + " not supported")
	}
}
