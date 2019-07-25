package graphs

import (
	"gmidarch/development/framework/messages"
)

type GraphExecutable struct {
	NumNodes        int
	EdgesExecutable [][]EdgeExecutable
}

type EdgeExecutableInfo struct {
	ActionName     string
	ActionChannel  *chan messages.SAMessage                                  // Channel
	Message        *messages.SAMessage                                       // Message
	ExternalAction func(*chan messages.SAMessage, *messages.SAMessage)       // External action
	InternalAction func(elem interface{}, name string, args ... interface{}) // Internal action
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
