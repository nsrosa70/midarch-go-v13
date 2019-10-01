package components

import (
	"gmidarch/development/framework/messages"
	"fmt"
	"os"
	"newsolution/development/element"
	"newsolution/development/miop"
	"newsolution/development/impl"
	"newsolution/shared/shared"
	"newsolution/development/artefacts/exec"
)

type Marshaller struct {
	CSP   string
	Graph exec.ExecGraph
}

func NewMarshaller() Marshaller {

	// create a new instance of Server
	r := new(Marshaller)

	return *r
}

func (m *Marshaller) Configure(invP, terP *chan messages.SAMessage) {

	// configure the state machine
	m.Graph = *exec.NewExecGraph(3)

	msg := new(messages.SAMessage)
	args := make([]*interface{}, 1)
	args[0] = new(interface{})
	*args[0] = msg

	newEdgeInfo := exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	m.Graph.AddEdge(0, 1, newEdgeInfo)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Process", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	m.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	m.Graph.AddEdge(2, 0, newEdgeInfo)

}

func (Marshaller) I_Process(msg *messages.SAMessage) {
	req := msg.Payload.(shared.Request)
	op := req.Op

	switch op {
	case "marshall":
		p1 := req.Args[0].(miop.Packet)
		r := impl.MarshallerImpl{}.Marshall(p1)
		*msg = messages.SAMessage{Payload: r}
	case "unmarshall":
		p1 := req.Args[0].([]byte)
		r := impl.MarshallerImpl{}.Unmarshall(p1)
		*msg = messages.SAMessage{Payload: r}
	default:
		fmt.Println("Marshaller:: Operation '" + op + "' not supported!!")
		os.Exit(0)
	}
}
