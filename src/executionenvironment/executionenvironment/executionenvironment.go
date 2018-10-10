package executionenvironment

import (
	"shared/conf"
	"shared/shared"
	"os"
	"verificationtools/fdr"
	"framework/configuration/configuration"
	"framework/message"
	"graph/execgraph"
	"reflect"
	"framework/library"
	"strings"
	"graph/fdrgraph"
	"fmt"
	"shared/errors"
	"framework/configuration/commands"
	"strconv"
	"shared/parameters"
	"framework/element"
)

type ExecutionEnvironment struct{}

func (ee ExecutionEnvironment) Deploy(confFile string) {

	// Load execution parameters
	shared.LoadParameters(os.Args[1:])

	// check of parameters
	shared.ShowExecutionParameters(false)

	// Generate Go configuration
	conf := conf.GenerateConf(confFile)

	// Initialize channels between Units and Adaptation manager
	channsUnits := make(map[string]chan commands.LowLevelCommand)
	for i := range conf.Components {
		id := conf.Components[i].Id
		channsUnits[id] = make(chan commands.LowLevelCommand)
	}
	for i := range conf.Connectors {
		id := conf.Connectors[i].Id
		channsUnits[id] = make(chan commands.LowLevelCommand)
	}

	// Initialize basic elements used throughout execution
	//channs := map[string]chan message.Message{}
	elemMaps := map[string]string{}

	// Configure channels & maps
	//channs = ee.ConfigureChannels(conf)
	elemMaps = ee.ConfigureMaps(conf)

	// Configure behaviour & behaviour expressions
	for i := range conf.Components {
		b := library.Repository[reflect.TypeOf(conf.Components[i].TypeElem).String()].CSP
		if b == "" {
			myError := errors.MyError{Source:"Execution Engine",Message:"Component '"+conf.Components[i].Id+"' does not exist in the Library"}
			myError.ERROR()
		}
		tempElem := element.Element{conf.Components[i].Id,element.Element{}.Behaviour,conf.Components[i].TypeElem,b}
		conf.Components[i] = tempElem
	}
	for i := range conf.Connectors {
		b := library.Repository[reflect.TypeOf(conf.Connectors[i].TypeElem).String()].CSP
		if b == "" {
			myError := errors.MyError{Source:"Execution Engine",Message:"Connector '"+conf.Connectors[i].Id+"'does not exist in the Library"}
			myError.ERROR()
		}
		tempElem := element.Element{conf.Connectors[i].Id,element.Element{}.Behaviour,conf.Connectors[i].TypeElem, b}
		conf.Connectors[i] = tempElem
	}

	// Check behaviour using FDR
	fdr := new(fdr.FDR)   // TODO
	ok := fdr.CheckBehaviour(conf,elemMaps)
	if !ok{
		myError := errors.MyError{Source:"Execution Engine",Message:"Configuration has a problem detected by FDR4"}
		myError.ERROR()
	}

	// Load graph generated by FDR (*.dot)
	fdrGraph := fdr.LoadFDRGraph(confFile)

	// Generate executable graph
	execGraph, execChannels := CreateExecGraph(fdrGraph)

	// Deploy configuration
	ee.Exec(conf, execChannels, execGraph)
}

func (ee ExecutionEnvironment) Exec(conf configuration.Configuration, channels map[string]chan message.Message, execGraph execgraph.GraphX) {

	// Start engine
	go StartEngine(execGraph)

	// Start components
	for e := range conf.Components {
		elemChannels := DefineChannels(channels, conf.Components[e].Id)
		actions := map[string][]string{}
		behaviour := library.Repository[reflect.TypeOf(conf.Components[e].TypeElem).String()].CSP
		actions[e] = FilterActions(strings.Split(behaviour, " "))
		individualChannels := map[string]chan message.Message{}
		for a := range actions[e] {
			individualChannels[actions[e][a]] = DefineChannel(elemChannels, actions[e][a])
		}
		go shared.Invoke(conf.Components[conf.Components[e].Id].TypeElem, "Loop", individualChannels)
	}
}

func FilterActions(actions []string) [] string {
	r := []string{}

	for a := range actions {
		action := actions[a]
		if strings.Contains(action, "I") || strings.Contains(action, "T") { // TODO
			if strings.Contains(action, ".") {
				action = action[:strings.Index(action, ".")]
			}
			r = append(r, action)
		}
	}
	return r
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
		fmt.Println("Error: channel '" + a + "' not found")
	}

	return r
}

func StartEngine(g execgraph.GraphX) {
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

func (ExecutionEnvironment) ConfigureChannels(conf configuration.Configuration) map[string]chan message.Message {
	channs := make(map[string]chan message.Message)

	for i := range conf.Attachments {
		c1Id := conf.Attachments[i].C1.Id
		c2Id := conf.Attachments[i].C2.Id
		tId := conf.Attachments[i].T.Id

		// c1 -> t
		key01 := c1Id + "." + "InvR" + "." + tId
		key02 := tId + "." + "InvP" + "." + c1Id
		key03 := tId + "." + "TerP" + "." + c1Id
		key04 := c1Id + "." + "TerR" + "." + tId
		channs[key01] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
		channs[key02] = channs[key01]
		channs[key03] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
		channs[key04] = channs[key03]

		// t -> c2
		key01 = tId + "." + "InvR" + "." + c2Id
		key02 = c2Id + "." + "InvP" + "." + tId
		key03 = c2Id + "." + "TerP" + "." + tId
		key04 = tId + "." + "TerR" + "." + c2Id
		channs[key01] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
		channs[key02] = channs[key01]
		channs[key03] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
		channs[key04] = channs[key03]
	}
	return channs
}

func (ExecutionEnvironment) ConfigureMaps(conf configuration.Configuration) (map[string]string) {

	elemMaps := make(map[string]string)
	partners := make(map[string]string)

	for i := range conf.Attachments {
		c1Id := conf.Attachments[i].C1.Id
		c2Id := conf.Attachments[i].C2.Id
		tId := conf.Attachments[i].T.Id
		if !strings.Contains(partners[c1Id], tId) {
			partners[c1Id] += ":" + tId
		}
		if !strings.Contains(partners[tId], c1Id) {
			partners[tId] += ":" + c1Id
		}
		if !strings.Contains(partners[tId], c2Id) {
			partners[tId] += ":" + c2Id
		}
		if !strings.Contains(partners[c2Id], tId) {
			partners[c2Id] += ":" + tId
		}
	}

	for i := range partners {
		p := strings.Split(partners[i], ":")
		c := 1
		for j := range p {
			if p[j] != "" {
				elemMaps[i+".e"+strconv.Itoa(c)] = p[j]
				c++
			}
		}
	}
	return elemMaps
}
