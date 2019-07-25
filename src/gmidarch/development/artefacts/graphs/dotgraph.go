package graphs

type GraphDot struct {
	NumNodes int
	EdgesDot    [][]EdgeDot
}

type EdgeDot struct {
	From   int
	To     int
	Action string
}

func NewGraphDot(n int) *GraphDot {
	return &GraphDot{
		NumNodes: n,
		EdgesDot:    make([][]EdgeDot, n),
	}
}

func (g *GraphDot) AddEdge(u, v int, a string) {
	g.EdgesDot[u] = append(g.EdgesDot[u], EdgeDot{From: u, To: v, Action: a})
}

func (g *GraphDot) AdjacentEdges(u int) []EdgeDot {
	return g.EdgesDot[u]
}
