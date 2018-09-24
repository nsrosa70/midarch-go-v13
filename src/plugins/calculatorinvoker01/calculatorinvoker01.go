package main

import (
	"framework/message"
	"apps/calculator/impl"
)

type CalculatorInvoker struct{}

var msg message.Message

func GetTypeElem() interface{}{
	return CalculatorInvoker{}
}

func GetBehaviourExp() string {
	//return library.BehaviourLibrary[calculatorinvoker.CalculatorInvoker{}]
	return "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"
}

func (CalculatorInvoker) I_PosInvP(msg *message.Message) {
	op := msg.Payload.(message.MIOP).RequestBody.Op
	args := msg.Payload.(message.MIOP).RequestBody.Args

	switch op {
	case "add":
		// prepare invocation
		argsX := args.([]interface{})
		p1 := int(argsX[0].(float64))
		p2 := int(argsX[1].(float64))

		//dispatch invocation
		//fmt.Println("[plugin] 1")
		r := impl.Add(p1, p2)

		//time.Sleep(5*time.Millisecond)

		// send reply
		header := message.ReplyHeader{1} // 1 - Success
		body := message.ReplyBody{Reply: r}
		miop := message.MIOP{ReplyHeader: header, ReplyBody: body}
		*msg = message.Message{miop}
	}
}