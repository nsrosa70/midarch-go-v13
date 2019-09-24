package components

import (
"gmidarch/development/artefacts/graphs"
"gmidarch/development/framework/messages"
	"newsolution/element"
)

type Executor struct {
	CSP      string
	Graph    graphs.GraphExecutable
	InvRChan chan messages.SAMessage
	Msg      messages.SAMessage
}

func NewExecutor(invR *chan messages.SAMessage) Executor {

	// create a new instance of client
	r := new(Executor)

	// configure the new instance
	r.CSP = "B = InvR -> B"
	r.InvRChan = *invR
	r.Msg = messages.SAMessage{}

	// configure the state machine
	r.Graph = *graphs.NewGraph(1)
	newEdgeInfo := graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: &r.InvRChan, Message: &r.Msg}
	r.Graph.AddEdge(0, 0, newEdgeInfo)

	return *r
}

