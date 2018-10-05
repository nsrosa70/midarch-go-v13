package execgraph

import "framework/message"

type Graph struct {
	NumNodes int
	Edges    [][]Edge
}

type GraphX struct {
	NumNodes int
	Edges    [][]EdgeX
}

type ExecAction struct{
	Action string
	Channel chan string
}

type ExecActionX struct{
	Action string
	Channel chan message.Message
}

type Edge struct {
	From   int
	To     int
	Action ExecAction
}

type EdgeX struct {
	From   int
	To     int
	Action ExecActionX
}

// NewGraph: Create graph with n nodes.
func NewGraph(n int) *Graph {
	return &Graph{
		NumNodes: n,
		Edges:    make([][]Edge, n),
	}
}

func NewGraphX(n int) *GraphX {
	return &GraphX{
		NumNodes: n,
		Edges:    make([][]EdgeX, n),
	}
}

func (g *Graph) AddEdge(u, v int, a ExecAction) {
	g.Edges[u] = append(g.Edges[u], Edge{From: u, To: v, Action: a})
}

func (g *GraphX) AddEdgeX(u, v int, a ExecActionX) {
	g.Edges[u] = append(g.Edges[u], EdgeX{From: u, To: v, Action: a})
}

func (g *Graph) AdjacentEdges(u int) []Edge {
	return g.Edges[u]
}

func (g *GraphX) AdjacentEdgesX(u int) []EdgeX {
	return g.Edges[u]
}