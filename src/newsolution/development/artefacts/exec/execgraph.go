package exec

import (
	"gmidarch/development/framework/messages"
)

type ExecGraph struct {
	NumNodes  int
	ExecEdges [][]ExecEdge
}

type ExecEdge struct {
	From int
	To   int
	Info ExecEdgeInfo
}

type TypeInternalAction func(any interface{}, name string, args [] *interface{})
type TypeExternalAction func(*chan messages.SAMessage, *messages.SAMessage)

type ExecEdgeInfo struct {
	ActionType     int // Internal & External
	ActionName     string
	ActionChannel  *chan messages.SAMessage // Channel
	Message        *messages.SAMessage      // Message
	ExternalAction TypeExternalAction       // External action
	InternalAction TypeInternalAction       // Internal action
	Args           [] *interface{}
	Response       *bool
}

func NewExecGraph(n int) *ExecGraph {
	return &ExecGraph{
		NumNodes:  n,
		ExecEdges: make([][]ExecEdge, n),
	}
}

func (g *ExecGraph) AddEdge(u, v int, a ExecEdgeInfo) {
	g.ExecEdges[u] = append(g.ExecEdges[u], ExecEdge{From: u, To: v, Info: a})
}

func (g *ExecGraph) AdjacentEdges(u int) []ExecEdge {
	return g.ExecEdges[u]
}
