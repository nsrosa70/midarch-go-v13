package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"newsolution/shared"
	"fmt"
	"os"
	"newsolution/element"
	"newsolution/impl"
)

type ServerCalculator struct {
	CSP      string
	Graph    graphs.GraphExecutable
	InvPChan chan messages.SAMessage
	TerPChan chan messages.SAMessage
	Msg      messages.SAMessage
}

func NewServerCalculator(invP *chan messages.SAMessage, terP *chan messages.SAMessage) ServerCalculator {

	// create a new instance of Server
	r := new(ServerCalculator)

	// configure the new instance
	r.CSP = "B = InvP -> I_Process -> TerP -> B"
	r.InvPChan = *invP
	r.TerPChan = *terP
	r.Msg = messages.SAMessage{}

	// configure the state machine
	r.Graph = *graphs.NewGraph(3)
	actionChannel := make(chan messages.SAMessage)

	newEdgeInfo := graphs.EdgeExecutableInfo{ExternalAction:element.Element{}.InvP, ActionType: 2, ActionChannel: &r.InvPChan, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: ServerCalculator{}.I_Process, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: &r.TerPChan, Message: &r.Msg}
	r.Graph.AddEdge(2, 0, newEdgeInfo)

	return *r
}

func (ServerCalculator) I_Process(msg *messages.SAMessage) {
	req := msg.Payload.(shared.Request)
	op := req.Op
	p1 := req.Args[0].(int)
	p2 := req.Args[1].(int)
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