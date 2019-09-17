package engine

import (
	"gmidarch/shared/shared"
	"fmt"
	"newsolution/components"
	"reflect"
)

type Engine struct {}

func (Engine) Execute(elem interface{}) {
	// Execute graph
	node := 0

	var e interface{}
	switch reflect.TypeOf(elem).String() {
	case "components.Client":
		e = elem.(components.Client)
	case "components.Unit":
		e = elem.(components.Unit)
	default
	}
	if (reflect.TypeOf(elem).String() == "components.Client"){

	} := elem.(components.Client)
	for {
		edges := e.Graph.AdjacentEdges(node)
		if len(edges) == 1 {
			edge := edges[0]
			if shared.IsInternal(edge.Info.ActionName) {
				edge.Info.InternalAction(e, edge.Info.ActionName, edge.Info.Message)
			} else {
				fmt.Println(edge.Info.ActionName)
			}
			node = edge.To
		} else {
			//msg := messages.SAMessage{}
			chosen := 0
			//choice(*elem, &msg, &chosen, edges)
			node = edges[chosen].To
		}
		if node == 0 {
			fmt.Println("Finished!!")
			break
		}
	}
	return
}