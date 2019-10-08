package connectors

import (
	"gmidarch/development/framework/messages"
	"newsolution/development/element"
	"newsolution/development/artefacts/exec"
)

type Requestreply struct {
	CSP      string
	Graph    exec.ExecGraph
}

func NewRequestReply() Requestreply {

	// create a new instance of client
	r := new(Requestreply)

	return *r
}

func (r *Requestreply) Configure (invP, terP, invR, terR *chan messages.SAMessage) Requestreply {

	// configure the new instance
	//msg := messages.SAMessage{}
	msg := new(messages.SAMessage)

	// configure the state machine
	r.Graph = *exec.NewExecGraph(4)
	newEdgeInfo := exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, Message: msg, ActionChannel: invP, ActionType: 2}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, Message: msg, ActionChannel: invR, ActionType: 2}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, Message: msg, ActionChannel: terR, ActionType: 2}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, Message: msg, ActionChannel: terP, ActionType: 2}
	r.Graph.AddEdge(3, 0, newEdgeInfo)

	return *r
}