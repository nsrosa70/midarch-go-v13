package element

import (
	"framework/messages"
	"shared/shared"
	"graph/execgraph"
	"graph/fdrgraph"
	"reflect"
)

type Element struct {
	Id        string
	TypeElem  interface{}
	CSP       string
	ExecGraph execgraph.Graph
	FDRGraph  fdrgraph.Graph
}

func (e *Element) SetFDRGraph(g fdrgraph.Graph) {
	e.FDRGraph = g
}

func (e *Element) SetExecGraph(g *execgraph.Graph) {
	e.ExecGraph = *g
}

func (Element) Loop(elem Element, graph execgraph.Graph) {

	// Execute graph
	node := 0
	for {
		edges := graph.AdjacentEdges(node)
		if len(edges) == 1 {
			if shared.IsInternal(edges[0].Action.ActionName) {
				edges[0].Action.InternalAction(elem.TypeElem, edges[0].Action.ActionName, edges[0].Action.Message)
			} else {
				edges[0].Action.ExternalAction(edges[0].Action.ActionChannel, edges[0].Action.Message)
			}
			node = edges[0].To
		} else {
			msg := messages.SAMessage{}
			chosen := 0
			choice(&msg, &chosen, edges)
			*edges[chosen].Action.Message = msg
			node = edges[chosen].To
		}
		if node == 0 {
			break
		}
	}
	return
}

func choice(msg *messages.SAMessage, chosen *int, edges []execgraph.Edge) {
	cases := make([]reflect.SelectCase, len(edges))

	var value reflect.Value
	for i := 0; i < len(edges); i++ {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Action.ActionChannel)}
	}
	*chosen, value, _ = reflect.Select(cases)
	*msg = value.Interface().(messages.SAMessage)

	cases = nil
}

// external actions common to all components and connectors
func (Element) InvP(invP *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*invP
}

func (Element) InvR(invR *chan messages.SAMessage, msg *messages.SAMessage) {
	*invR <- *msg
}

func (Element) TerR(terR *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*terR
}

func (Element) TerP(terP *chan messages.SAMessage, msg *messages.SAMessage) {
	*terP <- *msg
}