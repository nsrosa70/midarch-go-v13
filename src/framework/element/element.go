package element

import (
	"graph/execgraph"
	"graph/fdrgraph"
	"shared/shared"
	"framework/messages"
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
				r := false
				edges[0].Action.InternalAction(elem.TypeElem, edges[0].Action.ActionName, edges[0].Action.Message,&r)
			} else {
				edges[0].Action.ExternalAction(edges[0].Action.ActionChannel, edges[0].Action.Message)
			}
			node = edges[0].To
		} else {
			msg := messages.SAMessage{}
			chosen := 0
			choice(elem, &msg, &chosen, edges)
			node = edges[chosen].To
		}
		if node == 0 {
			break
		}
	}
	return
}

func choice(elem Element, msg *messages.SAMessage, chosen *int, edges []execgraph.Edge) {
	cases := make([]reflect.SelectCase, len(edges))

	// Assembly cases
	for i := 0; i < len(edges); i++ {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Action.ActionChannel)}
		if shared.IsInternal(edges[i].Action.ActionName) {

			// Execute Internal action
			r := false
			msgTemp := *edges[i].Action.Message
			edges[i].Action.InternalAction(elem.TypeElem, edges[i].Action.ActionName, &msgTemp, &r)

			// Update internal channel
			if r {
				go send(edges[i].Action.ActionChannel, msgTemp)
			}
		}
	}

	// Select channel (action)
	var value reflect.Value
	*chosen, value, _ = reflect.Select(cases)
	*edges[*chosen].Action.Message = value.Interface().(messages.SAMessage)
}

func send(channel *chan messages.SAMessage, msg messages.SAMessage) {
	*channel <- msg
}

func receive(channel *chan messages.SAMessage) {
	<-*channel
}

// external actions common to all components and connectors
func (elem Element) InvP(invP *chan messages.SAMessage, msg *messages.SAMessage) {
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
