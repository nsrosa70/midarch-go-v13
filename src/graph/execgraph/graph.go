package execgraph

type Graph struct {
	NumNodes int
	Edges    [][]Edge
}

type ExecAction struct{
	Action string
	Channel chan string
}

type Edge struct {
	From   int
	To     int
	Action ExecAction
}

// NewGraph: Create graph with n nodes.
func NewGraph(n int) *Graph {
	return &Graph{
		NumNodes: n,
		Edges:    make([][]Edge, n),
	}
}

func (g *Graph) AddEdge(u, v int, a ExecAction) {
	g.Edges[u] = append(g.Edges[u], Edge{From: u, To: v, Action: a})
}

func (g *Graph) AdjacentEdges(u int) []Edge {
	return g.Edges[u]
}