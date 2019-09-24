package components

import (
	"gmidarch/development/framework/messages"
	"gmidarch/development/artefacts/graphs"
	"newsolution/element"
	"newsolution/shared"
)

type CalculatorProxy struct {
	CSP   string
	Graph graphs.GraphExecutable
	InvP  chan messages.SAMessage
	TerP  chan messages.SAMessage
	InvR  chan messages.SAMessage
	TerR  chan messages.SAMessage
	Msg   messages.SAMessage
}

func NewCalculatorProxy(invP, terP, invR, terR *chan messages.SAMessage) CalculatorProxy {

	// create a new instance of Server
	r := new(CalculatorProxy)

	// configure the new instance
	r.CSP = "B = I_SetMessage -> I_Process -> B"
	r.InvP = *invP
	r.InvR = *invR
	r.TerR = *terR
	r.TerP = *terP
	r.Msg = messages.SAMessage{}

	// configure the state machine
	r.Graph = *graphs.NewGraph(6)
	actionChannel := make(chan messages.SAMessage)

	newEdgeInfo := graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: &r.InvP, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: CalculatorProxy{}.I_ProcessIn, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: &r.InvR, Message: &r.Msg}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: &r.TerR, Message: &r.Msg}
	r.Graph.AddEdge(3, 4, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: CalculatorProxy{}.I_ProcessOut, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(4, 5, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: &r.TerP, Message: &r.Msg}
	r.Graph.AddEdge(5, 0, newEdgeInfo)

	return *r
}

func (CalculatorProxy) I_ProcessIn(msg *messages.SAMessage) {
	inv := shared.Invocation{}
	inv.Host = "localhost"  // TODO
	inv.Port = 1313
	inv.Req = msg.Payload.(shared.Request)

	*msg = messages.SAMessage{Payload: inv}
}

func (CalculatorProxy) I_ProcessOut(msg *messages.SAMessage) {

	// check message

	//*msg = messages.SAMessage{Payload: msgFromServer}
}
