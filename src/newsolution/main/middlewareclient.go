package main

import (
	"fmt"
	"gmidarch/development/framework/messages"
	"newsolution/components"
	"newsolution/engine"
	"newsolution/shared"
	"newsolution/connectors"
	"sync"
)

func main() {
	var chn [17] chan messages.SAMessage

	for i := 0; i < 17; i++ {
		chn[i] = make(chan messages.SAMessage)
	}

	clientCalculator := components.NewClientCalculator(&chn[0], &chn[1])
	t1 := connectors.NewRequestReply(&chn[0], &chn[1], &chn[2], &chn[3])
	calculatorProxy := components.NewCalculatorProxy(&chn[2], &chn[3], &chn[4], &chn[5])

	t2 := connectors.NewRequestReply(&chn[4], &chn[5], &chn[6], &chn[7])
	requestor := components.NewRequestor(&chn[6], &chn[7], &chn[8], &chn[9], &chn[12], &chn[13])
	t3 := connectors.NewRequestReply(&chn[8], &chn[9], &chn[10], &chn[11])
	marshaller := components.NewMarshaller(&chn[10], &chn[11])
	t4 := connectors.NewRequestReply(&chn[12], &chn[13], &chn[14], &chn[15])
	crh := components.NewCRH(&chn[14], &chn[15])
	executor := components.NewExecutor(&chn[16])

	unit1 := components.NewUnit(clientCalculator, &chn[16])
	unit2 := components.NewUnit(calculatorProxy, &chn[16])
	unit3 := components.NewUnit(requestor, &chn[16])
	unit4 := components.NewUnit(marshaller, &chn[16])
	unit5 := components.NewUnit(crh, &chn[16])

	unit6 := components.NewUnit(t1, &chn[16])
	unit7 := components.NewUnit(t2, &chn[16])
	unit8 := components.NewUnit(t3, &chn[16])
	unit9 := components.NewUnit(t4, &chn[16])

	var wg sync.WaitGroup
	wg.Add(10)
	go engine.Engine{}.Execute(unit1, unit1.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit2, unit2.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit3, unit3.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit4, unit4.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit5, unit5.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit6, unit6.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit7, unit7.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit8, unit8.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit9, unit9.Graph, shared.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(executor, executor.Graph, shared.EXECUTE_FOREVER)
	wg.Wait()

	fmt.Scanln()
}
