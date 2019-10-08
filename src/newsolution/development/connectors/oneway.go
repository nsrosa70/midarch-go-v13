package connectors

import (
	"gmidarch/development/framework/messages"
	"newsolution/development/element"
	"newsolution/development/artefacts/exec"
)

type Oneway struct {
	CSP   string
	Graph exec.ExecGraph
}

func NewOneway() Oneway {

	// create a new instance of client
	r := new(Oneway)

	return *r
}

func (o *Oneway) ConfigureOneWay(invP, invR *chan messages.SAMessage) {

	// configure the state machine
	//msg := new(messages.SAMessage)
	msg := new(messages.SAMessage)
	args := make([]*interface{}, 1)
	args[0] = new(interface{})
	*args[0] = msg

	o.Graph = *exec.NewExecGraph(2)
	newEdgeInfo := exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, Message: msg, ActionChannel: invP, ActionType: 2}
	o.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, Message: msg, ActionChannel: invR, ActionType: 2}
	o.Graph.AddEdge(1, 0, newEdgeInfo)
}