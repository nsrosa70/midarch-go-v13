package main

import (
	"fmt"
	"gmidarch/development/framework/messages"
	"newsolution/development/components"
	"newsolution/development/connectors"
	"newsolution/execution/environment/engine"
	"newsolution/shared/parameters"
	"newsolution/execution/environment/frontend"
)

func oldmain() {
	frontend.FrontEnd{}.Deploy("senderreceiver.madl")

	fmt.Scanln()
}


func main() {
	chn1 := make(chan messages.SAMessage)
	chn2 := make(chan messages.SAMessage)
	chn3 := make(chan messages.SAMessage)

	sender := components.NewSender()
	sender.Configure(&chn1)
	t := connectors.NewOneway()
	t.ConfigureOneWay(&chn1, &chn2)
	receiver := components.NewReceiver()
	receiver.Configure(&chn2)
	executor := components.NewExecutor()
	executor.Configure(&chn3)

	unit1 := components.NewUnit()
	unit1.ConfigureUnit(sender,&chn3)
	unit2 := components.NewUnit()
	unit2.ConfigureUnit(receiver,&chn3)
	unit3 := components.NewUnit()
	unit3.ConfigureUnit(t,&chn3)

	go engine.Engine{}.Execute(unit1, unit1.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit2, unit2.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit3, unit3.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(executor, executor.Graph, parameters.EXECUTE_FOREVER)

	fmt.Scanln()
}
