package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
)

type Sender struct {
	CSP      string
	Graph    graphs.GraphExecutable
	InvRChan chan messages.SAMessage
	Msg      messages.SAMessage
}

func NewSender(chn *chan messages.SAMessage) Sender {

	// create a new instance of client
	r := new(Sender)

	// configure the new instance
	r.CSP = "B = I_SetMessage1 -> InvR -> B [] I_SetMessage2 -> InvR -> B [] I_SetMessage3 -> InvR -> B"
	r.InvRChan = *chn
	r.Msg = messages.SAMessage{}

	// configure the state machine
	r.Graph = *graphs.NewGraph(2)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo := graphs.EdgeExecutableInfo{InternalAction: Sender{}.I_SetMessage1, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: Sender{}.I_SetMessage2, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: Sender{}.I_SetMessage3, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: Sender{}.InvR, ActionType: 2, ActionChannel: &r.InvRChan, Message: &r.Msg}
	r.Graph.AddEdge(1, 0, newEdgeInfo)

	return *r
}

func (c Sender) I_SetMessage1(msg *messages.SAMessage) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 1)"}
}
func (c Sender) I_SetMessage2(msg *messages.SAMessage) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 2)"}
}
func (c Sender) I_SetMessage3(msg *messages.SAMessage) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 3)"}
}

func (c Sender) InvR(invR *chan messages.SAMessage, msg *messages.SAMessage) {
	*invR <- *msg
}
