package components

import (
	"gmidarch/development/framework/messages"
	"reflect"
	"os"
	"newsolution/development/element"
	"fmt"
	"newsolution/development/connectors"
	"newsolution/execution/environment/engine"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
	"newsolution/development/artefacts/exec"
)

type Unit struct {
	CSP   string
	Graph exec.ExecGraph
}

func NewUnit() Unit {

	// create a new instance of Unit
	r := new(Unit)

	return *r
}

func (u *Unit) ConfigureUnit(elem interface{}, invP *chan messages.SAMessage) {

	// configure the state machine
	u.Graph = *exec.NewExecGraph(2)
	msg := new(messages.SAMessage)
	actionChannel := make(chan messages.SAMessage)

	args := make([] *interface{}, 2)
	args[0] = new(interface{})
	args[1] = new(interface{})
	*args[0] = new(messages.SAMessage)
	*args[1] = elem

	newEdgeInfo := exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Execute", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	u.Graph.AddEdge(0, 0, newEdgeInfo)

	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionName: "InvP", ActionType: 2, ActionChannel: invP, Message: msg}
	u.Graph.AddEdge(0, 1, newEdgeInfo)

	actionChannel = make(chan messages.SAMessage)
	args1 := make([]*interface{}, 1)
	args1[0] = new(interface{})
	*args1[0] = new(messages.SAMessage)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_PerformAdaptation", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args1}
	u.Graph.AddEdge(1, 0, newEdgeInfo)
}

func (Unit) I_Execute(msg *messages.SAMessage, elem interface{}) {

	elemTemp := elem

	switch  reflect.TypeOf(elemTemp).String() {
	case "components.Sender":
		elemTemp := elemTemp.(Sender)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Receiver":
		elemTemp := elemTemp.(Receiver)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Client":
		elemTemp := elemTemp.(Client)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Server":
		elemTemp := elemTemp.(Server)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Calculatorclient":
		elemTemp := elemTemp.(Calculatorclient)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Calculatorserver":
		elemTemp := elemTemp.(Calculatorserver)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.CRH":
		elemTemp := elemTemp.(CRH)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.SRH":
		elemTemp := elemTemp.(SRH)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Marshaller":
		elemTemp := elemTemp.(Marshaller)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Requestor":
		elemTemp := elemTemp.(Requestor)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.CalculatorProxy":
		elemTemp := elemTemp.(CalculatorProxy)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Invoker":
		elemTemp := elemTemp.(Invoker)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)

	case "connectors.Oneway":
		elemTemp := elemTemp.(connectors.Oneway)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "connectors.Requestreply":
		elemTemp := elemTemp.(connectors.Requestreply)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	default:
		fmt.Println("Unit:: Architectural element '" + reflect.TypeOf(elemTemp).String() + "' not supported")
		os.Exit(0)
	}
}

func (Unit) I_PerformAdaptation(msg *messages.SAMessage) {
	fmt.Println("Unit:: Perform Adaptation ***********")
}
