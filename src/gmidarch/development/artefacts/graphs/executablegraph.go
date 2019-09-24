package graphs

import (
	"gmidarch/development/framework/messages"
)

type GraphExecutable struct {
	NumNodes        int
	EdgesExecutable [][]EdgeExecutable
}

type TypeInternalAction func(*messages.SAMessage)
type TypeInternalActionWithArgs func(*messages.SAMessage, [] *interface{})
type TypeExternalAction func(*chan messages.SAMessage, *messages.SAMessage)

type EdgeExecutableInfo struct {
	ActionType             int
	ActionName             string
	ActionChannel          *chan messages.SAMessage   // Channel
	Message                *messages.SAMessage        // Message
	ExternalAction         TypeExternalAction         // External action
	InternalAction         TypeInternalAction         // Internal action
	InternalActionWithArgs TypeInternalActionWithArgs // Internal action with args
	Args                   [] *interface{}
	Response               *bool
}

type EdgeExecutable struct {
	From int
	To   int
	Info EdgeExecutableInfo
}

func NewGraph(n int) *GraphExecutable {
	return &GraphExecutable{
		NumNodes:        n,
		EdgesExecutable: make([][]EdgeExecutable, n),
	}
}

func (g *GraphExecutable) AddEdge(u, v int, a EdgeExecutableInfo) {
	g.EdgesExecutable[u] = append(g.EdgesExecutable[u], EdgeExecutable{From: u, To: v, Info: a})
}

func (g *GraphExecutable) AdjacentEdges(u int) []EdgeExecutable {
	return g.EdgesExecutable[u]
}
