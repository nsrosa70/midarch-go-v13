package execgraph

import "framework/message"

type Graph struct {
	NumNodes int
	Edges    [][]Edge
}

type Action struct{
	ActionName string
	ActionChannel *chan message.Message
	Message *message.Message
	ExternalAction func(*chan message.Message, *message.Message)             // External action
	InternalAction func(elem interface{}, name string, args ... interface{}) // Internal action
}

type Edge struct {
	From   int
	To     int
	Action Action
}

func NewGraph(n int) *Graph {
	return &Graph{
		NumNodes: n,
		Edges:    make([][]Edge, n),
	}
}

func (g *Graph) AddEdge(u, v int, a Action) {
	g.Edges[u] = append(g.Edges[u], Edge{From: u, To: v, Action: a})
}

func (g *Graph) AdjacentEdges(u int) []Edge {
	return g.Edges[u]
}