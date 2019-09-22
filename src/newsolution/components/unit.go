package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"newsolution/engine"
	"reflect"
	"newsolution/connectors"
	"newsolution/shared"
	"fmt"
	"os"
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
	args := make([]interface{}, 1)
	args[0] = elem
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo := graphs.EdgeExecutableInfo{InternalActionWithArgs: Unit{}.I_Execute, ActionName: "I_EXECUTE", ActionType: 3, ActionChannel: &actionChannel, Message: &r.Msg, Args: args}
	r.Graph.AddEdge(0, 0, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: Unit{}.InvP, ActionType: 2, ActionName: "InvP", ActionChannel: &r.InvPChan, Message: &r.Msg}
	r.Graph.AddEdge(0, 0, newEdgeInfo)

	return *r
}

func (Unit) I_Execute(msg *messages.SAMessage, elem interface{}) {

	switch  reflect.TypeOf(elem.([]interface{})[0]).String() {
	case "components.Sender":
		elemTemp := elem.([]interface{})[0].(Sender)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.Receiver":
		elemTemp := elem.([]interface{})[0].(Receiver)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.Client":
		elemTemp := elem.([]interface{})[0].(Client)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "components.Server":
		elemTemp := elem.([]interface{})[0].(Server)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "connectors.OneWay":
		elemTemp := elem.([]interface{})[0].(connectors.OneWay)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	case "connectors.RequestReply":
		elemTemp := elem.([]interface{})[0].(connectors.RequestReply)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !shared.EXECUTE_FOREVER)
	default:
		fmt.Println("Unit:: Architectural element not supported")
		os.Exit(0)
	}
}

func (Unit) InvP(invP *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*invP
}
