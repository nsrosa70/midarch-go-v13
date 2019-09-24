package main

import (
	"fmt"
	"gmidarch/development/framework/messages"
	"newsolution/components"
	"newsolution/engine"
	"newsolution/shared"
	"newsolution/connectors"
)

func main() {
	var chn [13] chan messages.SAMessage

	for i := 0; i < 13; i++ {
		chn[i] = make(chan messages.SAMessage)
	}

	srh := components.NewSRH(&chn[0], &chn[1],"localhost",1313)
	t1 := connectors.NewRequestReply(&chn[0], &chn[1], &chn[2], &chn[3])
	invoker := components.NewInvoker(&chn[2],&chn[3],&chn[4],&chn[5],&chn[6],&chn[7])
	t2 := connectors.NewRequestReply(&chn[4], &chn[5], &chn[8], &chn[9])
	marshaller := components.NewMarshaller(&chn[8],&chn[9])
	t3 := connectors.NewRequestReply(&chn[6], &chn[7], &chn[10], &chn[11])
	calculatorServer := components.NewServerCalculator(&chn[10],&chn[11])
	executor := components.NewExecutor(&chn[12])

	unit1 := components.NewUnit(calculatorServer, &chn[12])
	unit2 := components.NewUnit(marshaller, &chn[12])
	unit3 := components.NewUnit(invoker, &chn[12])
	unit4 := components.NewUnit(srh, &chn[12])
	unit5 := components.NewUnit(t1, &chn[12])
	unit6 := components.NewUnit(t2, &chn[12])
	unit7 := components.NewUnit(t3, &chn[12])

	//var wg sync.WaitGroup
	//wg.Add(7)
	go engine.Engine{}.Execute(unit1, unit1.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit2, unit2.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit3, unit3.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit4, unit4.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit5, unit5.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit6, unit6.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit7, unit7.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(executor, executor.Graph, shared.EXECUTE_FOREVER)
	//wg.Wait()

	fmt.Scanln()
}
