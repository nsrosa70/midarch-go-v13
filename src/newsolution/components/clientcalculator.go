package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"fmt"
	"newsolution/shared"
	"newsolution/element"
)

type ClientCalculator struct {
	CSP      string
	Graph    graphs.GraphExecutable
	InvP chan messages.SAMessage
	TerP chan messages.SAMessage
	InvR chan messages.SAMessage
	TerR chan messages.SAMessage
	Msg      messages.SAMessage
}

func NewClientCalculator(invR, terR *chan messages.SAMessage) ClientCalculator {

	// create a new instance of client
	r := new(ClientCalculator)

	// configure the new instance
	r.CSP = "B = InvP -> I_ProcessIn -> InvR -> I_ProcessOut -> TerP -> B"
	r.InvR = *invR
	r.TerR = *terR
	r.Msg = messages.SAMessage{}

	// configure the state machine
	r.Graph = *graphs.NewGraph(5)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo := graphs.EdgeExecutableInfo{InternalAction: ClientCalculator{}.I_SetMessage, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: &r.InvR, Message: &r.Msg}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: &r.TerR, Message: &r.Msg}
	r.Graph.AddEdge(3, 4, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: ClientCalculator{}.I_PrintMessage, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(4, 0, newEdgeInfo)

	return *r
}

func (ClientCalculator) I_SetMessage(msg *messages.SAMessage) {
	args := make([]interface{},2)
	args[0] = 1
	args[1] = 2
	*msg = messages.SAMessage{Payload: shared.Request{Op:"add",Args:args}}
}

func (ClientCalculator) I_PrintMessage(msg *messages.SAMessage) {
	fmt.Println(msg.Payload)
}