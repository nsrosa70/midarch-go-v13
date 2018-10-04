package fdrgraph

type Graph struct {
	NumNodes int
	Edges    [][]Edge
}

type Edge struct {
	From   int
	To     int
	Action string
}

func NewGraph(n int) *Graph {
	return &Graph{
		NumNodes: n,
		Edges:    make([][]Edge, n),
	}
}

func (g *Graph) AddEdge(u, v int, a string) {
	g.Edges[u] = append(g.Edges[u], Edge{From: u, To: v, Action: a})
}

func (g *Graph) AdjacentEdges(u int) []Edge {
	return g.Edges[u]
}