package element

import (
	"gmidarch/development/framework/messages"
	"gmidarch/shared/shared"
	"reflect"
	"gmidarch/development/artefacts/graphs"
)

type ElementGo struct {
	ElemId string
	ElemType interface{}
	CSP string
	Info      interface{}
	GoStateMachine graphs.GraphExecutable
	FDRStateMachine graphs.GraphDot
}

// external actions common to all components and connectors
func (ElementGo) InvP(invP *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*invP
}

func (ElementGo) InvR(invR *chan messages.SAMessage, msg *messages.SAMessage) {
	*invR <- *msg
}

func (ElementGo) TerR(terR *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*terR
}

func (ElementGo) TerP(terP *chan messages.SAMessage, msg *messages.SAMessage) {
	*terP <- *msg
}

func (ElementGo) Loop(elem ElementGo, graph graphs.GraphExecutable) {

	// Execute graph
	node := 0

	for {
		edges := graph.AdjacentEdges(node)
		if len(edges) == 1 {
			edge := edges[0]
			if shared.IsInternal(edge.Info.ActionName) {
				r := false
				//fmt.Printf("%s:: %s \n",elem.ElemId,edge.Info.ActionName)
				edge.Info.InternalAction(elem.ElemType, edge.Info.ActionName, edge.Info.Message, &elem.Info, &r)
			} else {
				//fmt.Printf("%s:: %s %v \n",elem.ElemId,edge.Info.ActionName,edge.Info.ActionChannel)
				edge.Info.ExternalAction(edge.Info.ActionChannel, edge.Info.Message)
			}
			node = edge.To
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

func (e *ElementGo) SetInfo(info interface{}) {
	e.Info = info
}

func choice(elem ElementGo, msg *messages.SAMessage, chosen *int, edges []graphs.EdgeExecutable) {
	cases := make([]reflect.SelectCase, len(edges))

	// Assembly cases
	for i := 0; i < len(edges); i++ {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
		if shared.IsInternal(edges[i].Info.ActionName) {
			// Execute Internal action
			r := false
			edges[i].Info.InternalAction(elem.ElemType, edges[i].Info.ActionName, edges[i].Info.Message, &elem.Info, &r)

			// Update internal channel
			if r {
				go send(edges[i].Info.ActionChannel, *edges[i].Info.Message)
			}
		}
	}

	// Select channel (action)
	var value reflect.Value
	*chosen, value, _ = reflect.Select(cases)
	*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
}

func send(channel *chan messages.SAMessage, msg messages.SAMessage) {
	*channel <- msg
}