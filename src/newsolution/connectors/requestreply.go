package connectors

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"newsolution/element"
)

type RequestReply struct {
	CSP      string
	Graph    graphs.GraphExecutable
	Msg      messages.SAMessage
	InvP chan messages.SAMessage
	TerP chan messages.SAMessage
	InvR chan messages.SAMessage
	TerR chan messages.SAMessage
	Result bool
}

func NewRequestReply(invP, terP, invR, terR *chan messages.SAMessage) RequestReply {

	// create a new instance of client
	r := new(RequestReply)

	// configure the new instance
	r.CSP = "B = InvP -> InvR -> TerR -> TerP -> B"
	r.Msg = messages.SAMessage{Payload: ""}
	r.InvP= *invP
	r.InvR = *invR
	r.TerR = *terR
	r.TerP = *terP

	// configure the state machine
	r.Graph = *graphs.NewGraph(4)
	newEdgeInfo := graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvP, Message: &r.Msg, ActionChannel: &r.InvP, ActionType: 2}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvR, Message: &r.Msg, ActionChannel: &r.InvR, ActionType: 2}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerR, Message: &r.Msg, ActionChannel: &r.TerR, ActionType: 2}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerP, Message: &r.Msg, ActionChannel: &r.TerP, ActionType: 2}
	r.Graph.AddEdge(3, 0, newEdgeInfo)

	return *r
}