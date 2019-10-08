package components

import (
	"newsolution/development/artefacts/exec"
	"gmidarch/development/framework/messages"
	"newsolution/shared/shared"
	"fmt"
	"newsolution/development/element"
)

//var t1 time.Time
//var idx int

type Calculatorclient struct {
	CSP   string
	Graph exec.ExecGraph
}

func NewCalculatorclient() Calculatorclient {

	// create a new instance of client
	r := new(Calculatorclient)

	return *r
}

func (c *Calculatorclient) Configure(invR, terR *chan messages.SAMessage) {

	// configure the state machine
	c.Graph = *exec.NewExecGraph(4)
	actionChannel := make(chan messages.SAMessage)

	msg := new(messages.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	newEdgeInfo := exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Setmessage", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	c.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR, Message: msg}
	c.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR, Message: msg}
	c.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Printmessage", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	c.Graph.AddEdge(3, 0, newEdgeInfo)
}

func (Calculatorclient) I_Setmessage(msg *messages.SAMessage, info [] *interface{}) {

	argsTemp := make([]interface{}, 2)
	argsTemp[0] = 1
	argsTemp[1] = 2
	*msg = messages.SAMessage{Payload: shared.Request{Op: "add", Args: argsTemp}}

}

func (Calculatorclient) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {
	fmt.Println(msg.Payload)
}
