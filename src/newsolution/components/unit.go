package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"newsolution/engine"
	"reflect"
	"newsolution/connectors"
	"newsolution/shared"
	"os"
	"newsolution/element"
	"fmt"
)

type Unit struct {
	CSP      string
	Graph    graphs.GraphExecutable
	Element  interface{}
	InvPChan chan messages.SAMessage
	Msg      messages.SAMessage
}

func NewUnit(elem interface{}, invP *chan messages.SAMessage) Unit {

	// create a new instance of Unit
	r := new(Unit)

	// configure the new instance
	r.CSP = "B = I_Execute -> B [] InvP -> B"
	r.Msg = messages.SAMessage{}
	r.InvPChan = *invP

	// configure the state machine
	r.Graph = *graphs.NewGraph(1)
	args := make([] *interface{}, 1)
	args[0] = &elem
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo := graphs.EdgeExecutableInfo{InternalActionWithArgs: Unit{}.I_Execute, ActionName: "I_EXECUTE", ActionType: 3, ActionChannel: &actionChannel, Message: &r.Msg, Args: args}
	r.Graph.AddEdge(0, 0, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionName: "InvP", ActionChannel: &r.InvPChan, Message: &r.Msg}
	r.Graph.AddEdge(0, 0, newEdgeInfo)

	return *r
}

func (Unit) I_Execute(msg *messages.SAMessage, elem [] *interface{}) {

	elemTemp := *elem[0]
	switch  reflect.TypeOf(elemTemp).String() {
	case "components.Sender":
		//elemTemp := elem.([]interface{})[0].(Sender)
		elemTemp := elemTemp.(Sender)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.Receiver":
		//elemTemp := elem.([]interface{})[0].(Receiver)
		elemTemp := elemTemp.(Receiver)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.Client":
		//elemTemp := elem.([]interface{})[0].(Client)
		elemTemp := elemTemp.(Client)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.Server":
		//elemTemp := elem.([]interface{})[0].(Server)
		elemTemp := elemTemp.(Server)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.ClientCalculator":
		//elemTemp := elem.([]interface{})[0].(ClientCalculator)
		elemTemp := elemTemp.(ClientCalculator)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.ServerCalculator":
		//elemTemp := elem.([]interface{})[0].(ServerCalculator)
		elemTemp := elemTemp.(ServerCalculator)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.CRH":
		//elemTemp := elem.([]interface{})[0].(CRH)
		elemTemp := elemTemp.(CRH)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.SRH":
		//elemTemp := elem.([]interface{})[0].(SRH)
		elemTemp := elemTemp.(SRH)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.Marshaller":
		//elemTemp := elem.([]interface{})[0].(Marshaller)
		elemTemp := elemTemp.(Marshaller)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.Requestor":
		//elemTemp := elem.([]interface{})[0].(Requestor)
		elemTemp := elemTemp.(Requestor)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.CalculatorProxy":
		//elemTemp := elem.([]interface{})[0].(CalculatorProxy)
		elemTemp := elemTemp.(CalculatorProxy)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.Invoker":
		//elemTemp := elem.([]interface{})[0].(CalculatorProxy)
		elemTemp := elemTemp.(Invoker)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)

	case "connectors.OneWay":
		//elemTemp := elem.([]interface{})[0].(connectors.OneWay)
		elemTemp := elemTemp.(connectors.OneWay)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "connectors.RequestReply":
		//elemTemp := elem.([]interface{})[0].(connectors.RequestReply)
		elemTemp := elemTemp.(connectors.RequestReply)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	default:
		//fmt.Println("Unit:: Architectural element '"+reflect.TypeOf(elem.([]interface{})[0]).String()+"' not supported")
		fmt.Println("Unit:: Architectural element '"+reflect.TypeOf(elemTemp).String()+"' not supported")
		os.Exit(0)
	}
}