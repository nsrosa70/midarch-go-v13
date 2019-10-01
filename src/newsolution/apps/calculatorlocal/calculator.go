package calculatorlocal

import (
	"fmt"
	"gmidarch/development/framework/messages"
	"newsolution/development/components"
	"newsolution/engine"
	"newsolution/development/connectors"
	"newsolution/shared"
)

func main() {
	chn1 := make(chan messages.SAMessage)
	chn2 := make(chan messages.SAMessage)
	chn3 := make(chan messages.SAMessage)
	chn4 := make(chan messages.SAMessage)
	chn5 := make(chan messages.SAMessage)

	client := components.NewClientCalculator(&chn1, &chn2)
	t := connectors.NewRequestReply(&chn1, &chn3, &chn4, &chn2)
	server := components.NewServerCalculator(&chn3, &chn4)
	executor := components.NewExecutor(&chn5)

	unit1 := components.NewUnit(client, &chn5)
	unit2 := components.NewUnit(server, &chn5)
	unit3 := components.NewUnit(t, &chn5)

	go engine.Engine{}.Execute(unit1, unit1.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit2, unit2.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit3, unit3.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(executor, executor.Graph, shared.EXECUTE_FOREVER)

	//go engine.Engine{}.Execute(sender,sender.Graph)
	//go engine.Engine{}.Execute(receiver,receiver.Graph)
	//go engine.Engine{}.Execute(t,t.Graph)

	fmt.Scanln()
}
