package wgraph

import (
//	"shared/shared"
	"shared/shared"
)

type Graph struct {
	NumNodes int
	Edges    [][]Edge
}

type Edge struct {
	From   int
	To     int
	Action shared.ParMapActions
}

// NewGraph: Create graph with n nodes.
func NewGraph(n int) *Graph {
	return &Graph{
		NumNodes: n,
		Edges:    make([][]Edge, n),
	}
}

func (g *Graph) AddEdge(u, v int, a shared.ParMapActions) {
	g.Edges[u] = append(g.Edges[u], Edge{From: u, To: v, Action: a})
}

func (g *Graph) AdjacentEdges(u int) []Edge {
	return g.Edges[u]
}

type GraphSimple struct {
	NumNodes int
	Edges    [][]EdgeSimple
}

type EdgeSimple struct {
	From   int
	To     int
	Action string
}

// NewGraph: Create graph with n nodes.
func NewGraphSimple(n int) *GraphSimple {
	return &GraphSimple{
		NumNodes: n,
		Edges:    make([][]EdgeSimple, n),
	}
}

func (g *GraphSimple) AddEdgeSimple(u, v int, a string) {
	g.Edges[u] = append(g.Edges[u], EdgeSimple{From: u, To: v, Action: a})
}

func (g *GraphSimple) AdjacentEdgesSimple(u int) []EdgeSimple {
	return g.Edges[u]
}