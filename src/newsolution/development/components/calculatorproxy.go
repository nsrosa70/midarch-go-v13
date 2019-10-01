package components

import (
	"gmidarch/development/framework/messages"
	"newsolution/development/element"
	"newsolution/shared/shared"
	"newsolution/development/artefacts/exec"
)

type CalculatorProxy struct {
	CSP   string
	Graph exec.ExecGraph
}

func NewCalculatorProxy() CalculatorProxy {

	// create a new instance of Server
	r := new(CalculatorProxy)

	return *r
}

func (c *CalculatorProxy) Configure(invP, terP, invR, terR *chan messages.SAMessage) CalculatorProxy {

	// configure the state machine
	c.Graph = *exec.NewExecGraph(6)
	actionChannel := make(chan messages.SAMessage)
	msg := new(messages.SAMessage)
	args := make([]*interface{}, 1)
	args[0] = new(interface{})
	*args[0] = msg

	newEdgeInfo := exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	c.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName:"I_ProcessIn",ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	c.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR, Message: msg}
	c.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR, Message: msg}
	c.Graph.AddEdge(3, 4, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName:"I_ProcessOut",ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	c.Graph.AddEdge(4, 5, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	c.Graph.AddEdge(5, 0, newEdgeInfo)

	return *c
}

func (CalculatorProxy) I_ProcessIn(msg *messages.SAMessage) {
	inv := shared.Invocation{}
	inv.Host = "localhost"  // TODO
	inv.Port = 1313
	inv.Req = msg.Payload.(shared.Request)

	*msg = messages.SAMessage{Payload: inv}
}

func (CalculatorProxy) I_ProcessOut(msg *messages.SAMessage) {

	result := msg.Payload.([]interface{})
	*msg = messages.SAMessage{Payload: int(result[0].(float64))}
}
