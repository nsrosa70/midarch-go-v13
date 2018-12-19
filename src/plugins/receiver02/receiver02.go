package main

import (
"framework/messages"
"fmt"
)

type Receiver struct{}

func GetTypeElem() interface{}{
	return Receiver{}
}

func GetBehaviourExp() string {
	//return libraries.BehaviourLibrary[calculatorinvoker.CalculatorInvoker{}]
	return "B = InvP.e1 -> I_PosInvP -> B"
}

func (Receiver) I_PosInvP(msg *messages.SAMessage, info interface{}, r *bool) {
	//fmt.Println("Receiver:: " + msg.Payload.(string))
	fmt.Printf("Receiver [TYPE 02]:: %v\n",msg.Payload)
	return
}

