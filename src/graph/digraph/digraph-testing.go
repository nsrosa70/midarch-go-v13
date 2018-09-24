package digraph1

import (
	"testing"
	"fmt"
)

func setupTest() (g *Graph, u, v *Node) {
	g = NewGraph()
	u = NewNode("Wikipedia")
	v = NewNode("Article")
	return
}

func TestBasicOperation(t *testing.T) {
	g,u,v := setupTest()
	var ok bool

	fmt.Println("Test: Add Single Node")
	g.AddNodes(u, NewNode(""))
	if ok, _=g.Has(u); !ok {
		t.Errorf("Expected to have %s [node] in graph. ", u.ToString())
	}

	o := NewNode("Wikipedia")
	if ok, _=g.Has(u); !ok {
		t.Errorf("Expected to have %s [node] already in graph. ", o.ToString())
	}

	check := g.EdgeExist(u, NewNode(""))
	if check == true {
		t.Errorf("Expected not to have edge exist for `%s` and nil object", u.ToString())
	}

	check = g.EdgeExist(v, u)
	if check == true {
		t.Errorf("Expected not to have edge exist for `%s` and `%s`", v.ToString(), u.ToString())
	}

	check = g.EdgeExist(u, v)
	if check == true {
		t.Errorf("Expected not to have edge exist for `%s` and `%s`", u.ToString(), v.ToString())
	}

	if _,l:=g.Nodes(); l > 1 {
		t.Errorf("Expected single node to be in the graph.")
	}

	fmt.Println("Test: Add Nodes with edge.")
	g.AddNodes(u, v)
	if ok, _ =g.Has(u); !ok {
		t.Errorf("Expected to have %s [node] in graph. ", u.ToString())
	}

	if ok, _ =g.Has(u); !ok {
		t.Errorf("Expected to have %s [node] in graph. ", v.ToString())
	}

	if !g.EdgeExist(u, v) {
		t.Errorf("Expected to have Edge exist between `%s` and `%s`", u.ToString(), v.ToString())
	}

	x := NewNode("Outlier")
	if ok, _ = g.Has(x); ok {
		t.Errorf("Expected not to have `%s` node.", x.ToString())
	}

	if g.EdgeExist(u, x) {
		t.Errorf("Expected not to have Edge exist between `%s` and `%s`", u.ToString(), x.ToString())
	}

	fmt.Println("Test: single node connected with multiple nodes")
	g.AddNodes(u, x)
	if _, l := g.Edges(u); l != 2 {
		t.Errorf("Expected to have two nodes connected with %s", u.ToString())
	}

	fmt.Println("Test: Chaining of node. x<-u->v->x")
	g.AddNodes(v, x)
	if _, l := g.Nodes(); l!=3 {
		t.Errorf("Expected to have two nodes in a graph")
	}
}

func TestConcurrentSafety(t *testing.T) {
	g := NewGraph()
	u := NewNode("one")
	v := NewNode("hundred")
	y := NewNode("thousand")
	go func() {
		g.AddNodes(u, v)
		g.AddNodes(u, y)
	}()

	for {
		if ok, _ := g.Has(u); ok {
			g.Edges(u)
		}
	}
}

func BenchmarkAddSingleNode(b *testing.B) {
	g, u, _ := setupTest()
	for i:=0; i<b.N; i++ {
		g.AddNodes(u, NewNode(""))
	}
}

func BenchmarkAddNodeWithEdge(b *testing.B) {
	g, u, v := setupTest()
	for i:=0; i<b.N; i++ {
		g.AddNodes(u, v)
	}
}

func BenchmarkNodeConnectedWithTwoNode(b *testing.B) {
	g, u, v := setupTest()
	x := NewNode("X")
	for i:=0; i<b.N; i++ {
		g.AddNodes(u, v)
		g.AddNodes(u, x)
	}
}

func BenchmarkAddChainOfNodes(b *testing.B) {
	g, u, v := setupTest()
	x := NewNode("Chaining")
	for i:=0; i<b.N; i++ {
		g.AddNodes(u,v)
		g.AddNodes(v,x)
	}
}
