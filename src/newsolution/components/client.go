package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"fmt"
	"newsolution/element"
)

type Client struct {
	CSP      string
	Graph    graphs.GraphExecutable
	InvRChan chan messages.SAMessage
	TerRChan chan messages.SAMessage
	Msg      messages.SAMessage
}

func NewClient(invR *chan messages.SAMessage, terR *chan messages.SAMessage) Client {

	// create a new instance of client
	r := new(Client)

	// configure the new instance
	r.CSP = "B = I_SetMessage1 -> InvR -> TerR -> B"
	r.InvRChan = *invR
	r.TerRChan = *terR
	r.Msg = messages.SAMessage{}

	// configure the state machine
	r.Graph = *graphs.NewGraph(4)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo := graphs.EdgeExecutableInfo{InternalAction: Client{}.I_SetMessage1, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: &r.InvRChan, Message: &r.Msg}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: &r.TerRChan, Message: &r.Msg}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: Client{}.I_PrintMessage, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(3, 0, newEdgeInfo)

	return *r
}

func (Client) I_SetMessage1(msg *messages.SAMessage) {
	*msg = messages.SAMessage{Payload: "Hello World (Type 1)"}
}
func (Client) I_PrintMessage(msg *messages.SAMessage) {
	fmt.Println(msg.Payload)
}