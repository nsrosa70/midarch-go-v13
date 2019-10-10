package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"newsolution/shared/shared"
	"newsolution/gmidarch/development/element"
	"fmt"
)

type Client struct {
	CSP   string
	Graph graphs.ExecGraph
}

func NewClient() Client {

	// create a new instance of client
	r := new(Client)

	return *r
}

func (c *Client) Configure(invR, terR *chan messages.SAMessage) {

	// configure the state machine
	c.Graph = *graphs.NewExecGraph(4)
	actionChannel := make(chan messages.SAMessage)

	msg := new(messages.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	newEdgeInfo := graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Setmessage", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	c.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionName: "InvR", ActionType: 2, ActionChannel: invR, Message:msg}
	c.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionName: "TerR", ActionType: 2, ActionChannel: terR, Message:msg}
	c.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Printmessage", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	c.Graph.AddEdge(3, 0, newEdgeInfo)
}

func (Client) I_Setmessage(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: "Hello World from Client"}
}

func (Client) I_Printmessage(msg *messages.SAMessage, info [] *interface{}) {

	fmt.Println(msg.Payload)
}