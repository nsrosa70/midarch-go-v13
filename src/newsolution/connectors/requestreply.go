package connectors

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
)

type RequestReply struct {
	CSP      string
	Graph    graphs.GraphExecutable
	Msg      messages.SAMessage
	InvPChan chan messages.SAMessage
	TerPChan chan messages.SAMessage
	InvRChan chan messages.SAMessage
	TerRChan chan messages.SAMessage
	Result bool
}

func NewRequestReply(invP *chan messages.SAMessage, invR *chan messages.SAMessage, terR *chan messages.SAMessage, terP *chan messages.SAMessage) RequestReply {

	// create a new instance of client
	r := new(RequestReply)

	// configure the new instance
	r.CSP = "B = InvP -> InvR -> TerR -> TerP -> B"
	r.Msg = messages.SAMessage{Payload: ""}
	r.InvPChan = *invP
	r.InvRChan = *invR
	r.TerRChan = *terR
	r.TerPChan = *terP

	// configure the state machine
	r.Graph = *graphs.NewGraph(4)
	newEdgeInfo := graphs.EdgeExecutableInfo{ExternalAction: RequestReply{}.InvP, Message: &r.Msg, ActionChannel: &r.InvPChan, ActionType: 2}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: RequestReply{}.InvR, Message: &r.Msg, ActionChannel: &r.InvRChan, ActionType: 2}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: RequestReply{}.TerR, Message: &r.Msg, ActionChannel: &r.TerRChan, ActionType: 2}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: RequestReply{}.TerP, Message: &r.Msg, ActionChannel: &r.TerPChan, ActionType: 2}
	r.Graph.AddEdge(3, 0, newEdgeInfo)

	return *r
}

func (c RequestReply) InvP(invP *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*invP
}

func (c RequestReply) InvR(invR *chan messages.SAMessage, msg *messages.SAMessage) {
	*invR <- *msg
}

func (c RequestReply) TerR(terR *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*terR
}

func (c RequestReply) TerP(terP *chan messages.SAMessage, msg *messages.SAMessage) {
	*terP <- *msg
}
