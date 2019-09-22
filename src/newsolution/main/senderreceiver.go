package main

import (
	"fmt"
	"gmidarch/development/framework/messages"
	"newsolution/components"
	"newsolution/engine"
	"newsolution/connectors"
	"newsolution/shared"
)

func main() {
	chn1 := make(chan messages.SAMessage)
	chn2 := make(chan messages.SAMessage)
	chn3 := make(chan messages.SAMessage)

	sender := components.NewSender(&chn1)
	receiver := components.NewReceiver(&chn2)
	executor := components.NewExecutor(&chn3)
	t := connectors.NewOneWay(&chn1, &chn2)

	unit1 := components.NewUnit(sender,&chn3)
	unit2 := components.NewUnit(receiver,&chn3)
	unit3 := components.NewUnit(t,&chn3)

	go engine.Engine{}.Execute(unit1, unit1.Graph,shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit2, unit2.Graph,shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit3, unit3.Graph,shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(executor,executor.Graph,shared.EXECUTE_FOREVER)

	//go engine.Engine{}.Execute(sender,sender.Graph)
	//go engine.Engine{}.Execute(receiver,receiver.Graph)
	//go engine.Engine{}.Execute(t,t.Graph)

	fmt.Scanln()
}
