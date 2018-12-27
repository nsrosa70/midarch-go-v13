package execgraph

import "framework/messages"

type Graph struct {
	NumNodes int
	Edges    [][]Edge
}

type EdgeInfo struct {
	ActionName     string
	ActionChannel  *chan messages.SAMessage
	Message        *messages.SAMessage
	ExternalAction func(*chan messages.SAMessage, *messages.SAMessage) // External action
	//InternalAction func(elem interface{}, name string, args ... interface{}) // Internal action
	InternalAction func(elem interface{}, name string, args ... interface{}) // Internal action
}

type Edge struct {
	From int
	To   int
	Info EdgeInfo
}

func NewGraph(n int) *Graph {
	return &Graph{
		NumNodes: n,
		Edges:    make([][]Edge, n),
	}
}

func (g *Graph) AddEdge(u, v int, a EdgeInfo) {
	g.Edges[u] = append(g.Edges[u], Edge{From: u, To: v, Info: a})
}

func (g *Graph) AdjacentEdges(u int) []Edge {
	return g.Edges[u]
}
