package shared

import (
	"reflect"
	"framework/message"
	"strings"
	"os"
	"time"
	"shared/parameters"
	"fmt"
	"strconv"
	"graph/fdrgraph"
	"graph/execgraph"
)

const PREFIX_ACTION = "->"
const CHOICE = "[]"

type Invocation struct {
	Method  reflect.Value
	InArgs  []reflect.Value
	OutArgs [] reflect.Value
}

type MonitoredCorrectiveData string  // used in channel Monitor -> Analyser (Corrective)
type MonitoredEvolutiveData []string // used in channel Monitor -> Analyser (Evolutive)

type AnalysisResult struct {
	// used in channel Analyser -> Planner
	Result   interface{}
	Analysis int
}

var ValidActions = map[string]bool{
	"InvP": true,
	"TerP": true,
	"InvR": true,
	"TerR": true}

func IsInternal(action string) bool {
	return action[0:2] == "I_"
}

func ToActions(behaviour string) [] string {
	// B = InvP -> B [] InvP -> B
	var actions []string

	if !strings.Contains(behaviour, "[]") {
		behaviourTemp := strings.Split(behaviour, "=")
		behaviour = behaviourTemp[1][0:strings.LastIndex(behaviourTemp[1], "->")]
		actions = strings.Split(behaviour, PREFIX_ACTION)
		if len(actions) == 0 {
			actions[0] = strings.TrimSpace(behaviour)
		} else {
			for i := range actions {
				actions[i] = strings.TrimSpace(actions[i])
			}
		}
	} else {
		behaviourTemp := strings.Split(behaviour, "=")
		branches := strings.Split(behaviourTemp[1], "[]")
		idx := 0
		for i := range branches {
			actionsTemp := strings.Split(branches[i], PREFIX_ACTION)
			for j := range actionsTemp {
				action := strings.TrimSpace(actionsTemp[j])
				if action != "B" && action != "" {
					actions = append(actions, strings.TrimSpace(action))
					idx++
				}
			}
		}
	}
	return actions
}

type ParMapActions struct {
	F1 func(*chan message.Message, *message.Message)            // External action
	F2 func(any interface{}, name string, args ... interface{}) // Internal action
	P1 *message.Message
	P2 *chan message.Message
	P3 interface{}
	P4 string
	P5 *chan string // REMOVE
	P6 *string      // REMOVE
}

type SubMessage struct {
	I int
}
type MyMessage struct {
	Name string
	Age  int
	X    interface{}
}

func Invoke(any interface{}, name string, args ... interface{}) {
	inputs := make([]reflect.Value, len(args))

	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)

	inputs = nil
	return
}

func SelectChan(action string, id string, channs map[string]chan message.Message, elemMaps map[string]string) chan message.Message {

	p1 := action[0:strings.Index(action, ".")]
	p2 := action[strings.Index(action, ".")+1 : len(action)]

	keyMaps := id + "." + p2
	keyChannel := id + "." + p1 + "." + elemMaps[keyMaps]

	return channs[keyChannel]
}

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func LoadParameters(args []string){
	for i:= range args{
		variable := strings.Split(args[i],"=")
		switch strings.TrimSpace(variable[0]) {
		case "SAMPLE_SIZE":
			parameters.SAMPLE_SIZE,_ = strconv.Atoi(variable[1])
		case "REQUEST_TIME":
			temp1,_ := strconv.Atoi(variable[1])
			temp2   := time.Duration(temp1)
			parameters.REQUEST_TIME = temp2
		case "INJECTION_TIME":
			temp1,_ := strconv.Atoi(variable[1])
			temp2   := time.Duration(temp1)
			parameters.INJECTION_TIME = temp2
		case "MONITOR_TIME":
			temp1,_ := strconv.Atoi(variable[1])
			temp2   := time.Duration(temp1)
			parameters.MONITOR_TIME = temp2
		case "STRATEGY":
			parameters.STRATEGY,_ = strconv.Atoi(variable[1])
		case "IS_ADAPTIVE":
			parameters.IS_ADAPTIVE,_ = strconv.ParseBool(variable[1])
		case "NAMING_HOST":
			parameters.NAMING_HOST = variable[1]
		case "QUEUEING_HOST":
			parameters.QUEUEING_HOST = variable[1]
		default:
			fmt.Println("Shared:: Parameter '"+variable[0]+"' does not exist")
			os.Exit(0)
		}
	}
}

func CreateExecGraph(fdrGraph fdrgraph.Graph) (execgraph.GraphX, map[string]chan message.Message) {
	graph := execgraph.NewGraphX(fdrGraph.NumNodes)
	channels := map[string]chan message.Message{}

	// create channels
	for e1 := range fdrGraph.Edges {
		for e2 := range fdrGraph.Edges[e1] {
			eTemp := fdrGraph.Edges[e1][e2]
			if _, ok := channels[eTemp.Action]; !ok {
				channels[eTemp.Action] = make(chan message.Message)
			}
			graph.AddEdgeX(eTemp.From, eTemp.To, execgraph.ExecActionX{Action: eTemp.Action, Channel: channels[eTemp.Action]})
		}
	}
	return *graph, channels
}

func DefineChannels(channels map[string]chan message.Message, elem string) map[string]chan message.Message {
	r := map[string]chan message.Message{}

	for c := range channels {
		if strings.Contains(c, elem) {
			r[c] = channels[c]
		}
	}
	return r
}

func DefineChannel(channels map[string]chan message.Message, a string) chan message.Message {
	var r chan message.Message
	found := false

	for c := range channels {
		if (a[:2] != "I_") {
			if strings.Contains(c, a) && c[:2] != "I_" {
				r = channels[c]
				found = true
				break
			}
		} else {
			if strings.Contains(c, a) {
				r = channels[c]
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println("Error: channel '" + a + " not found")
	}

	return r
}

func Control(g execgraph.GraphX) {
	node := 0
	var msg = message.Message{}
	for {
		edges := g.AdjacentEdgesX(node)
		if len(edges) == 1 { // one edge
			node = edges[0].To
			if IsToElement(edges[0].Action.Action) {
				edges[0].Action.Channel <- msg
			} else {
				msg = <-edges[0].Action.Channel
			}
		} else { // two+ edges
			chosen := 0
			ChoiceX(&msg, &chosen, edges)
			node = edges[chosen].To
		}
	}
}

func IsToElement(action string) bool {
	if action[:2] == "I_" || action[:4] == "InvP" || action[:4] == "TerR" {
		return true
	} else { // TerP and InvR
		return false
	}
}

func ChoiceX(msg *message.Message, chosen *int, edges []execgraph.EdgeX) {
	cases := make([]reflect.SelectCase, len(edges))
	var value reflect.Value

	for i := 0; i < len(edges); i++ {
		if IsToElement(edges[i].Action.Action) {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(edges[i].Action.Channel), Send: reflect.ValueOf(*msg)}
		} else {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(edges[i].Action.Channel), Send: reflect.Value{}}
		}
	}

	*chosen, value, _ = reflect.Select(cases)
	if !IsToElement(edges[*chosen].Action.Action) {
		*msg = value.Interface().(message.Message)
	}
	cases = nil
}



