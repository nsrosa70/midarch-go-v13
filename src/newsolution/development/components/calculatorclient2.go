package components

import (
	"gmidarch/development/framework/messages"
	"newsolution/development/element"
	"fmt"
	"newsolution/shared/shared"
	"newsolution/development/artefacts/exec"
)

//var t1 time.Time
//var idx int

type Calculatorclient2 struct {
	CSP   string
	Graph exec.ExecGraph
}

func NewClientCalculator2(invR, terR *chan messages.SAMessage) Calculatorclient2 {

	// create a new instance of client
	r := new(Calculatorclient2)
	msg := new(messages.SAMessage)

	// configure the state machine
	r.Graph = *exec.NewExecGraph(4)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo := exec.ExecEdgeInfo{InternalAction: shared.Invoke,ActionName:"I_SetMessage", ActionType: 1, ActionChannel: &actionChannel, Message: msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR, Message: msg}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR, Message: msg}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke,ActionName:"I_PrintMessage", ActionType: 1, ActionChannel: &actionChannel, Message: msg}
	r.Graph.AddEdge(3, 0, newEdgeInfo)

	return *r
}

func (Calculatorclient2) I_SetMessage(msg *messages.SAMessage, args [] *interface{}) {

//	if idx < 1 {
//		idx++
		//t1 = time.Now()
		argsTemp := make([]interface{}, 2)
		argsTemp[0] = 13
		argsTemp[1] = 13
		*msg = messages.SAMessage{Payload: shared.Request{Op: "add", Args: argsTemp}}
//	} else {
//		os.Exit(0)
//	}
}

func (Calculatorclient2) I_PrintMessage(msg *messages.SAMessage, args [] *interface{}) {

	//fmt.Println(time.Now().Sub(t1))
	fmt.Println(msg.Payload)
}
