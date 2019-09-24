package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"newsolution/shared"
	"fmt"
	"os"
	"newsolution/element"
	"newsolution/miop"
	"newsolution/impl"
)

type Marshaller struct {
	CSP   string
	Graph graphs.GraphExecutable
	InvP  chan messages.SAMessage
	TerP  chan messages.SAMessage
	Msg   messages.SAMessage
}

func NewMarshaller(invP, terP *chan messages.SAMessage) Marshaller {

	// create a new instance of Server
	r := new(Marshaller)

	// configure the new instance
	r.CSP = "B = InvP -> I_Process -> TerP -> B"
	r.InvP = *invP
	r.TerP = *terP
	r.Msg = messages.SAMessage{}

	// configure the state machine
	r.Graph = *graphs.NewGraph(3)

	newEdgeInfo := graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: &r.InvP, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: Marshaller{}.I_Process, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: &r.TerP, Message: &r.Msg}
	r.Graph.AddEdge(2, 0, newEdgeInfo)

	return *r
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
