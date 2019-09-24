package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"strings"
	"newsolution/element"
)

type Server struct {
	CSP      string
	Graph    graphs.GraphExecutable
	InvPChan chan messages.SAMessage
	TerPChan chan messages.SAMessage
	Msg      messages.SAMessage
}

func NewServer(invP *chan messages.SAMessage, terP *chan messages.SAMessage) Server {

	// create a new instance of Server
	r := new(Server)

	// configure the new instance
	r.CSP = "B = I_SetMessage1 -> InvR -> TerR -> B"
	r.InvPChan = *invP
	r.TerPChan = *terP
	r.Msg = messages.SAMessage{}

	// configure the state machine
	r.Graph = *graphs.NewGraph(3)
	actionChannel := make(chan messages.SAMessage)

	newEdgeInfo := graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: &r.InvPChan, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: Server{}.I_Process, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: &r.TerPChan, Message: &r.Msg}
	r.Graph.AddEdge(2, 0, newEdgeInfo)

	return *r
}

func (Server) I_Process(msg *messages.SAMessage) {
	msgTemp := strings.ToUpper(msg.Payload.(string))
	*msg = messages.SAMessage{Payload: msgTemp}
}
