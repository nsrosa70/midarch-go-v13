package element

import (
	"framework/message"
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

	// execute graph
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
			msg := message.Message{}
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

func choice(msg *message.Message, chosen *int, edges []execgraph.Edge) {
	cases := make([]reflect.SelectCase, len(edges))

	var value reflect.Value
	for i := 0; i < len(edges); i++ {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Action.ActionChannel)}
	}
	*chosen, value, _ = reflect.Select(cases)
	*msg = value.Interface().(message.Message)

	cases = nil
}

// external actions common to all components and connectors
func (Element) InvP(invP *chan message.Message, msg *message.Message) {
	*msg = <-*invP
}

func (Element) InvR(invR *chan message.Message, msg *message.Message) {
	*invR <- *msg
}

func (Element) TerR(terR *chan message.Message, msg *message.Message) {
	*msg = <-*terR
}

func (Element) TerP(terP *chan message.Message, msg *message.Message) {
	*terP <- *msg
}