package engine

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"reflect"
	"newsolution/shared"
)

type Engine struct{}

func (Engine) Execute(elem interface{}, graph graphs.GraphExecutable, executionMode bool) {

	// Execute graph
	node := 0

	for {
		edges := graph.AdjacentEdges(node)
		if len(edges) == 1 {
			edge := edges[0]
			switch edge.Info.ActionType {
			case 1:
				edge.Info.InternalAction(edge.Info.Message)
			case 2:
				edge.Info.ExternalAction(edge.Info.ActionChannel, edge.Info.Message)
			case 3:
				edge.Info.InternalActionWithArgs(edge.Info.Message, edge.Info.Args, )
			}
			node = edge.To
		} else {
			msg := messages.SAMessage{}
			chosen := 0
			choice(elem, &msg, &chosen, edges)
			node = edges[chosen].To
		}
		if node == 0 && executionMode != shared.EXECUTE_FOREVER {
			break
		}
	}
	return
}

func choice(elem interface{}, msg *messages.SAMessage, chosen *int, edges []graphs.EdgeExecutable) {
	casesExternal := make([]reflect.SelectCase, len(edges)+1)
	casesInternal := make([]reflect.SelectCase, len(edges))

	// Assembly cases
	for i := 0; i < len(edges); i++ {
		switch edges[i].Info.ActionType {
		case 1: // Internal
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
			edges[i].Info.InternalAction(edges[i].Info.Message)
			go send(edges[i].Info.ActionChannel, *edges[i].Info.Message)

		case 2: // External
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
		case 3: // Internal with Arguments

			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
			edges[i].Info.InternalActionWithArgs(edges[i].Info.Message, edges[i].Info.Args)
			go send(edges[i].Info.ActionChannel, *edges[i].Info.Message)
		}
	}

	// default case
	casesExternal[len(edges)] = reflect.SelectCase{Dir: reflect.SelectDefault}

	// select external first
	var value reflect.Value
	*chosen, value, _ = reflect.Select(casesExternal)

	if *chosen != (len(edges)) { // NOT DEFAULT
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
		return
	} else {
		*chosen, value, _ = reflect.Select(casesInternal)
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	}
}

/*
func choiceOriginal(elem interface{}, msg *messages.SAMessage, chosen *int, edges []graphs.EdgeExecutable) {
	casesExternal := make([]reflect.SelectCase, len(edges)+1)
	casesInternal := make([]reflect.SelectCase, len(edges))

	// Assembly cases
	for i := 0; i < len(edges); i++ {
		switch edges[i].Info.ActionType {
		case 1: // Internal
			r := false
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
			edges[i].Info.InternalAction(edges[i].Info.Message)
			if r {
				go send(edges[i].Info.ActionChannel, *edges[i].Info.Message)
			}
		case 2: // External
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
		case 3: // Internal with Arguments
			r := false
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
			edges[i].Info.InternalActionWithArgs(edges[i].Info.Message, &r, edges[i].Info.Args)
			if r {
				go send(edges[i].Info.ActionChannel, *edges[i].Info.Message)
			}
		}
	}

	// default case
	casesExternal[len(edges)] = reflect.SelectCase{Dir: reflect.SelectDefault}

	// select external first
	var value reflect.Value
	*chosen, value, _ = reflect.Select(casesExternal)

	if *chosen != (len(edges)) { // NOT DEFAULT
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
		return
	} else {
		*chosen, value, _ = reflect.Select(casesInternal)
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	}
}
*/
func send(channel *chan messages.SAMessage, msg messages.SAMessage) {
	*channel <- msg
}
