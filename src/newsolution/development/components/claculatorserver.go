package components

import (
	"gmidarch/development/framework/messages"
	"fmt"
	"os"
	"newsolution/development/element"
	"newsolution/development/impl"
	"newsolution/shared/shared"
	"newsolution/development/artefacts/exec"
)

type Calculatorserver struct {
	CSP      string
	Graph    exec.ExecGraph
}

func Newcalculatorserver() Calculatorserver {

	// create a new instance of Server
	r := new(Calculatorserver)

	return *r
}

func (c *Calculatorserver) Configure(invP *chan messages.SAMessage, terP *chan messages.SAMessage) {

	// configure the state machine
	c.Graph = *exec.NewExecGraph(3)

	msg := new(messages.SAMessage)
	args := make([]*interface{}, 1)
	args[0] = new(interface{})
	*args[0] = msg

	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo := exec.ExecEdgeInfo{ExternalAction:element.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	c.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke,ActionName:"I_Process", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	c.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	c.Graph.AddEdge(2, 0, newEdgeInfo)
}

func (Calculatorserver) I_Process(msg *messages.SAMessage) {
	req := msg.Payload.(shared.Request)
	op := req.Op
	p1 := int(req.Args[0].(float64))
	p2 := int(req.Args[1].(float64))
	r := 0

	switch op {
	case "add":
		r = impl.CalculatorImpl{}.Add(p1,p2)
	case "sub":
		r = impl.CalculatorImpl{}.Sub(p1,p2)
	case "mul":
		r = impl.CalculatorImpl{}.Mul(p1,p2)
	case "div":
		r = impl.CalculatorImpl{}.Div(p1,p2)
	default:
		fmt.Println("Server Calculator:: Operation '" + op + "' not supported!!")
		os.Exit(0)
	}

	*msg = messages.SAMessage{Payload: r}
}