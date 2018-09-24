package executionunit

import (
	"framework/message"
	"framework/element"
	"fmt"
	"strings"
	"reflect"
	"framework/library"
	"graph/wgraph"
	"shared/shared"
	"shared/errors"
	"framework/configuration/commands"
	"shared/parameters"
)

var msg message.Message

type ExecutionUnit struct{}

var numReplaces = 0

func (ExecutionUnit) Exec(elem element.Element, channs map[string]chan message.Message, elemMaps map[string]string, chanUnit chan commands.LowLevelCommand) {
	elem.Behaviour = elem.BehaviourExec
	elem.BehaviourExp = library.BehaviourLibrary[reflect.TypeOf(elem.TypeElem).String()]

	// check behaviour
	id := elem.Id
	if !isBehaviourValid(elem) {
		myError := errors.MyError{Source: "Execution Unit", Message: "Behaviour of '" + id + "' is not VALID\n"}
		myError.ERROR()
	}

	// generate graph
	graph := wgraph.Graph{}
	GenerateGraph(&graph, elem.BehaviourExp, elem, msg, channs, elemMaps)

	// execute behaviour
	for {
		shared.Invoke(elem, "BehaviourExec", elem, graph, msg, channs, elemMaps)
		select {
		case cmd := <-chanUnit: // a new management action is received
			switch cmd.Cmd {
			case commands.REPLACE_COMPONENT:
				//fmt.Println("Unit:: Replace")
				elem = cmd.Args
			case commands.STOP:
				fmt.Println("Unit:: STOP [TODO]")
			}
		default: // no new management action received
		}
	}
}

func isBehaviourValid(elem element.Element) bool {
	e := true

	// empty behaviours
	if elem.BehaviourExp == "" {
		e = false
		return e
	}
	actions := shared.ToActions(elem.BehaviourExp)
	if len(actions) == 0 {
		e = false
		return e
	}

	// valid external actions
	for i := range actions {
		if !shared.IsInternal(actions[i]) {
			actionTemp := actions[i][0:strings.Index(actions[i], ".")]
			if !shared.ValidActions[actionTemp] {
				e = false
				return e
			}
		}
	}

	// valid internal  actions
	elemType := reflect.TypeOf(elem.TypeElem)
	for i := range actions {
		if shared.IsInternal(actions[i]) {
			foundMethod := false
			for j := 0; j < elemType.NumMethod(); j++ {
				method := elemType.Method(j)
				if (method.Name == actions[i]) {
					foundMethod = true
					break
				}
			}
			if !foundMethod {
				myError := errors.MyError{Source: "Execution Unit", Message: "Unit:: Internal Method '" + actions[i] + "' of '" + elem.Id + "' behaviour does not exist in its interface"}
				myError.ERROR()
			}
		}
	}

	return e
}

func ProcessBranch(g *wgraph.Graph, b string, nextTo *int, elem element.Element, msg message.Message, channs map[string]chan message.Message, elemMaps map[string]string) {
	actions := strings.Split(b, "->")
	from := 0
	to := 0
	eActions := shared.ParMapActions{}
	for i := 0; i < len(actions)-1; i++ {
		action := strings.TrimSpace(actions[i])
		nextAction := strings.TrimSpace(actions[i+1])
		if nextAction == "B" {
			from = to
			to = 0
		} else {
			from = to
			to = *nextTo + from + 1
		}
		*nextTo = from
		id := elem.Id
		if shared.IsInternal(action) { // using reflection
			params := shared.ParMapActions{F2: shared.Invoke, P3: elem.TypeElem, P4: action, P1: &msg}
			mapType := shared.ParMapActions{}
			mapType = params
			eActions = mapType
		} else {
			params := shared.ParMapActions{}
			actionTemp := action[0:strings.Index(action, ".")]
			switch actionTemp {
			case "InvR":
				invr := shared.SelectChan(action, id, channs, elemMaps)
				params = shared.ParMapActions{F1: element.Element{}.InvR, P1: &msg, P2: &invr, P4: action}
			case "TerR":
				terr := shared.SelectChan(action, id, channs, elemMaps)
				params = shared.ParMapActions{F1: new(element.Element).TerR, P1: &msg, P2: &terr, P4: action}
			case "InvP":
				invp := shared.SelectChan(action, id, channs, elemMaps)
				params = shared.ParMapActions{F1: new(element.Element).InvP, P1: &msg, P2: &invp, P4: action}
			case "TerP":
				terp := shared.SelectChan(action, id, channs, elemMaps)
				params = shared.ParMapActions{F1: new(element.Element).TerP, P1: &msg, P2: &terp, P4: action}
			}
			mapType := shared.ParMapActions{}
			mapType = params
			eActions = mapType
		}
		g.AddEdge(from, to, eActions)
	}
}

func ProcessBranchSimple(g *wgraph.GraphSimple, b string, nextTo *int) {
	actions := strings.Split(b, "->")
	from := 0
	to := 0
	for i := 0; i < len(actions)-1; i++ {
		action := strings.TrimSpace(actions[i])
		nextAction := strings.TrimSpace(actions[i+1])
		if nextAction == "B" {
			from = to
			to = 0
		} else {
			from = to
			to = *nextTo + from + 1
		}
		g.AddEdgeSimple(from, to, action)
		*nextTo = from
	}
}

func GenerateGraph(g *wgraph.Graph, b string, elem element.Element, msg message.Message, channs map[string]chan message.Message, elemMaps map[string]string) {
	// Only behaviours like B = B1 [] B2 [] .... [] Bn
	tempGraph := wgraph.NewGraph(parameters.GRAPH_SIZE) // TODO

	b = b[strings.Index(b, "=")+1:]
	nextTo := 0
	if strings.Contains(b, "[]") {
		branches := strings.Split(b, "[]")
		for i := range branches {
			ProcessBranch(tempGraph, branches[i], &nextTo, elem, msg, channs, elemMaps)
		}
	} else {
		ProcessBranch(tempGraph, b, &nextTo, elem, msg, channs, elemMaps)
	}
	*g = *tempGraph
}

func GenerateGraphSimple(g *wgraph.GraphSimple, b string) {
	// Only behaviours like B = B1 [] B2 [] .... [] Bn
	tempGraph := wgraph.NewGraphSimple(parameters.GRAPH_SIZE)

	b = b[strings.Index(b, "=")+1:]
	nextTo := 0
	if strings.Contains(b, "[]") {
		branches := strings.Split(b, "[]")
		for i := range branches {
			ProcessBranchSimple(tempGraph, branches[i], &nextTo)
		}
	} else {
		ProcessBranchSimple(tempGraph, b, &nextTo)
	}
	*g = *tempGraph
}
