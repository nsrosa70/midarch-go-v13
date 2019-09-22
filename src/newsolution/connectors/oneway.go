package connectors

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
)

type OneWay struct {
	CSP      string
	Graph    graphs.GraphExecutable
	Msg      messages.SAMessage
	InvPChan chan messages.SAMessage
	InvRChan chan messages.SAMessage
}

func NewOneWay(invP *chan messages.SAMessage, invR *chan messages.SAMessage) OneWay {

	// create a new instance of client
	r := new(OneWay)

	// configure the new instance
	r.CSP = "B = InvP -> InvR -> B"
	r.Msg = messages.SAMessage{Payload: ""}
	r.InvPChan = *invP
	r.InvRChan = *invR

	// configure the state machine
	r.Graph = *graphs.NewGraph(2)
	newEdgeInfo := graphs.EdgeExecutableInfo{ExternalAction: OneWay{}.InvP, Message: &r.Msg, ActionChannel: &r.InvPChan, ActionType: 2}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: OneWay{}.InvR, Message: &r.Msg, ActionChannel: &r.InvRChan, ActionType: 2}
	r.Graph.AddEdge(1, 0, newEdgeInfo)

	return *r
}

func (c OneWay) InvP(invP *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*invP
}

func (c OneWay) InvR(invR *chan messages.SAMessage, msg *messages.SAMessage) {
	*invR <- *msg
}
