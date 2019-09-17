package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
)

type Client struct {
	CSP string
	Graph graphs.GraphExecutable
	Buffer messages.SAMessage
	InvRChan chan messages.SAMessage
	TerRChan chan messages.SAMessage
}

func NewClient() Client{
	r := new(Client)

	r.CSP = "B = InvR -> TerR -> B"
	r.Buffer = messages.SAMessage{Payload:""}
	r.InvRChan = make(chan messages.SAMessage,1)
	r.TerRChan = make(chan messages.SAMessage,1)

	r.Graph = *graphs.NewGraph(2)
	newEdgeInfo := graphs.EdgeExecutableInfo{ActionName:"InvR"}
	r.Graph.AddEdge(0,1,newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ActionName:"TerR"}
	r.Graph.AddEdge(1,0,newEdgeInfo)
	return *r
}

func (c *Client) I_PreInvR(){
	c.Buffer = messages.SAMessage{Payload:"Hello World"}
}

func (c *Client) InvR(){
	c.InvRChan <- c.Buffer
}

func (c *Client) TerR(){
	c.Buffer = <- c.TerRChan
}
