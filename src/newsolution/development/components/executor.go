package components

import (
	"newsolution/development/artefacts/exec"
	"time"
	"fmt"
	"gmidarch/development/framework/messages"
	"newsolution/shared/shared"
	"newsolution/development/element"
)

type Executor struct {
	CSP      string
	Graph    exec.ExecGraph
}

func NewExecutor() Executor {

	// create a new instance of client
	r := new(Executor)

	return *r
}


func (e *Executor) Configure(invR *chan messages.SAMessage) {

	// configure the state machine
	e.Graph = *exec.NewExecGraph(2)
	actionChannel := make(chan messages.SAMessage)

	msg := new(messages.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	newEdgeInfo := exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Newadaptation", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info:info}
	e.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR, Message: msg}
	e.Graph.AddEdge(1, 0, newEdgeInfo)
}

func (Executor) I_Newadaptation(msg *messages.SAMessage) {
	time.Sleep(1000000 * time.Millisecond)
	fmt.Println("Executor:: NEW ADAPTATION AVAILABLE ******************")
	*msg = *msg
}
