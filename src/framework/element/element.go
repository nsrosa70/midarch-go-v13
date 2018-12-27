package element

import (
	"graph/execgraph"
	"graph/fdrgraph"
	"shared/shared"
	"framework/messages"
	"reflect"
	"strings"
	"shared/errors"
)

type Element struct {
	Id        string
	TypeElem  interface{}
	CSP       string
	FDRGraph  fdrgraph.Graph
	ExecGraph execgraph.Graph
	Info      interface{} // Any particular information necessary for the proper functioning of the component
}

func (e *Element) SetInfo(info interface{}) {
	e.Info = info
}

func (e *Element) SetFDRGraph(g fdrgraph.Graph) {
	e.FDRGraph = g
}

func (e *Element) SetExecGraph(g *execgraph.Graph) {
	e.ExecGraph = *g
}

func (Element) Loop(elem Element, graph *execgraph.Graph) {

	// Execute graph
	node := 0

	for {
		edges := graph.AdjacentEdges(node)
		if len(edges) == 1 {
			edge := edges[0]
			if shared.IsInternal(edge.Info.ActionName) {
				r := false
				//fmt.Printf("%v start %v %v\n",time.Now(),elem.Id,edges[0].Action.ActionName)
				//if reflect.TypeOf(elem.TypeElem).String() == "components.Requestor" {
				//	fmt.Printf("Element:: [BEFORE][%v] [%v] [%v] [%p] \n", elem.Id, edge.Info.ActionName, edge.Info.Message, edge.Info.Message)
				//}
				//edge.Info.InternalAction(elem.TypeElem, edge.Info.ActionName, edge.Info.Message, &elem.Info, &r)
				edge.Info.InternalAction(elem.TypeElem, edge.Info.ActionName, edge.Info.Message, &elem.Info, &r)
				//if reflect.TypeOf(elem.TypeElem).String() == "components.Requestor" {
				//	fmt.Printf("Element:: [AFTER][%v] [%v] [%v] [%p] \n", elem.Id, edge.Info.ActionName, *edge.Info.Message, edge.Info.Message)
				//}
				//fmt.Printf("%v complete %v %v\n",time.Now(),elem.Id,edges[0].Action.ActionName)
			} else {
				//if reflect.TypeOf(elem.TypeElem).String() == "components.Requestor" {
				//	fmt.Printf("Element:: [BEFORE][%v] [%v] [%v] [%p] \n", elem.Id, edge.Info.ActionName, edge.Info.Message, edge.Info.Message)
				//}
				edge.Info.ExternalAction(edge.Info.ActionChannel, edge.Info.Message)
				//if reflect.TypeOf(elem.TypeElem).String() == "components.Requestor" {
				//	fmt.Printf("Element:: [AFTER][%v] [%v] [%v] [%p] \n", elem.Id, edge.Info.ActionName, edge.Info.Message, edge.Info.Message)
				//}
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

func choice(elem Element, msg *messages.SAMessage, chosen *int, edges []execgraph.Edge) {
	cases := make([]reflect.SelectCase, len(edges))

	// Assembly cases
	for i := 0; i < len(edges); i++ {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
		if shared.IsInternal(edges[i].Info.ActionName) {
			// Execute Internal action
			r := false
			edges[i].Info.InternalAction(elem.TypeElem, edges[i].Info.ActionName, edges[i].Info.Message, &elem.Info, &r)

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

func DefineOldElement(comp interface{}, newElement interface{}) string {
	found := false
	oldElementId := ""
	components := comp.(map[string]Element)

	// TODO check compatibility of old and new elements by type
	for i := range components {
		oldElementType := reflect.TypeOf(components[i].TypeElem).String()
		oldElementType = oldElementType[strings.LastIndex(oldElementType, ".")+1 : len(oldElementType)]
		newElementType := reflect.TypeOf(newElement).String()
		newElementType = newElementType[strings.LastIndex(newElementType, ".")+1 : len(newElementType)]
		if oldElementType == newElementType {
			oldElementId = components[i].Id
			found = true
		}
	}

	if !found {
		myError := errors.MyError{Source: "Planner", Message: "New and old components have not COMPATIBLE types"}
		myError.ERROR()
	}
	return oldElementId
}
