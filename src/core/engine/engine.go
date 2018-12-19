package engine

import (
	"shared/shared"
	"os"
	"verificationtools/fdr"
	"framework/configuration/configuration"
	"framework/messages"
	"graph/execgraph"
	"strings"
	"shared/errors"
	"framework/configuration/commands"
	"strconv"
	"shared/parameters"
	"framework/element"
	"fmt"
	"reflect"
	"framework/libraries"
	"core/versioninginjector"
)

type Engine struct{}

func (engine Engine) Deploy(adlFileName string) {

	// Initialize environment
	engine.Initialization()

	// Prepare configuration to be executed
	appConf := engine.PrepareConfiguration(adlFileName)
	eeConf := engine.PrepareConfiguration("ExecutionEnvironment.conf")

	// Start App Configuration
	engine.ConfigureExecutionEnvironmentInfo(&eeConf, appConf) // TODO - Improve
	engine.ConfigureElementInfo(&appConf)                      // TODO - Improve
	engine.ConfigureUnits(&eeConf, appConf)
	engine.StartConfiguration(eeConf)
	go versioninginjector.InjectAdaptiveEvolution(parameters.PLUGIN_BASE_NAME)
}

func (Engine) ConfigureUnits(eeConf *configuration.Configuration, appConf configuration.Configuration) {
	availableUnits := []string{}

	// Identify units
	for u := range eeConf.Components {
		if reflect.TypeOf(eeConf.Components[u].TypeElem).String() == "components.ExecutionUnit" { // TODO Improve
			availableUnits = append(availableUnits, eeConf.Components[u].Id)
		}
	}

	if len(availableUnits) != len(appConf.Components)+len(appConf.Connectors) {
		fmt.Println("Engine:: Available units are not enough to execute the components/connectos. Hence, there is a problem in the configuration")
		os.Exit(0)
	}

	// Associate info to units, e.g., u1.info = c1, u2.info = t, u3.info = c2
	idx := 0
	for c := range appConf.Components {
		elem := eeConf.Components[availableUnits[idx]]
		elem.SetInfo(appConf.Components[c])
		eeConf.Components[availableUnits[idx]] = elem
		idx++
	}

	for t := range appConf.Connectors {
		elem := eeConf.Components[availableUnits[idx]]
		elem.SetInfo(appConf.Connectors[t])
		eeConf.Components[availableUnits[idx]] = elem
		idx++
	}
}

func (Engine) ConfigureExecutionEnvironmentInfo(eeConf *configuration.Configuration, appConf configuration.Configuration) {

	// Only components of the execution environment are checked (no connectors)
	for c1 := range eeConf.Components {
		switch reflect.TypeOf(eeConf.Components[c1].TypeElem).String() {
		case "components.ExecutionEnvironment": // TODO improve
			listOfElements := []element.Element{}
			for c2 := range appConf.Components {
				listOfElements = append(listOfElements, appConf.Components[c2])
			}
			for t := range appConf.Connectors {
				listOfElements = append(listOfElements, appConf.Connectors[t])
			}
			elem := eeConf.Components[c1]
			elem.SetInfo(listOfElements)
			eeConf.Components[c1] = elem
		case "components.MAPEKPlanner":
			elem := eeConf.Components[c1]
			elem.SetInfo(appConf.Components)  // only components are replaceable
			eeConf.Components[c1] = elem
		default:
			elem := eeConf.Components[c1]
			elem.SetInfo("none")
			eeConf.Components[c1] = elem
		}
	}
}

func (Engine) ConfigureElementInfo(conf *configuration.Configuration) {

	// Components only
	for c := range conf.Components {
		elem := conf.Components[c]
		elem.SetInfo("none")
		conf.Components[c] = elem
	}
}

func StartElement(elem element.Element) {
	for {
		shared.Invoke(elem, "Loop", elem, elem.ExecGraph)
	}
}

func (Engine) StartConfiguration(conf configuration.Configuration) {
	for c := range conf.Components {
		msg := messages.SAMessage{}
		for e1 := range conf.Components[c].ExecGraph.Edges {
			for e2 := range conf.Components[c].ExecGraph.Edges[e1] {
				conf.Components[c].ExecGraph.Edges[e1][e2].Action.Message = &msg
			}
		}
		go StartElement(conf.Components[c])
	}

	for t := range conf.Connectors {
		msg := messages.SAMessage{}
		for e1 := range conf.Connectors[t].ExecGraph.Edges {
			for e2 := range conf.Connectors[t].ExecGraph.Edges[e1] {
				conf.Connectors[t].ExecGraph.Edges[e1][e2].Action.Message = &msg
			}
		}
		go StartElement(conf.Connectors[t])
	}
}

/*
func (environment ExecutionEnvironment) StartConfigurationOld(conf configuration.Configuration, managementChannels map[string]chan commands.LowLevelCommand) {
	// Start execution units
	for c := range conf.Components {
		go executionunit.ExecutionUnit{}.Exec(conf.Components[c], managementChannels[conf.Components[c].Id])
	}
	for t := range conf.Connectors {
		go executionunit.ExecutionUnit{}.Exec(conf.Connectors[t], managementChannels[conf.Connectors[t].Id])
	}
}
*/
func (engine Engine) PrepareConfiguration(adlFileName string) configuration.Configuration {

	// Generate Go configuration
	conf := configuration.MapADLIntoGo(adlFileName)

	// Configure structural channels and maps of components/connectors
	engine.ConfigureStructuralChannelsAndMaps(&conf)

	// Check behaviour using FDR
	fdrChecker := new(fdr.FDR)
	ok := fdrChecker.CheckBehaviour(conf)
	if !ok {
		myError := errors.MyError{Source: "Execution Engine", Message: "Configuration has a problem detected by FDR4"}
		myError.ERROR()
	}

	// Generate *.dot files
	// FDR.GenerateFDRGraphs()  // TODO

	// Load graph generated by FDR (*.dot)
	fdrChecker.LoadFDRGraphs(&conf)

	// Generate executable graph
	engine.CreateExecGraphs(&conf)

	// Check if actions and their respective implementations exist
	CheckActionsAndImplementations(conf)

	return conf

}

func (ee Engine) Initialization() {
	// Load execution parameters
	shared.LoadParameters(os.Args[1:])

	// Perform checks on the library of c
	libraries.CheckLibrary()

	// Show execution parameters
	shared.ShowExecutionParameters(false)
}

func CheckActionsAndImplementations(conf configuration.Configuration) {

	// Check components
	for c := range conf.Components {
		for e1 := range conf.Components[c].ExecGraph.Edges {
			for e2 := range conf.Components[c].ExecGraph.Edges[e1] {
				action := conf.Components[c].ExecGraph.Edges[e1][e2].Action.ActionName
				if shared.IsExternal(action) {
					if action != shared.INVP && action != shared.TERP && action != shared.INVR && action != shared.TERR {
						fmt.Println("EE:: Component '" + conf.Components[c].Id + "' has an invalid action: '" + action)
						os.Exit(0)
					}
				} else {
					if shared.IsInternal(action) {
						st := reflect.TypeOf(conf.Components[c].TypeElem)
						_, ok := st.MethodByName(action)
						if !ok {
							fmt.Println("EE: Component '" + conf.Components[c].Id + "' has an invalid action: '" + action + "'")
							os.Exit(0)
						}

					} else {
						fmt.Println("EE: Component '" + conf.Components[c].Id + "' has an invalid action: '" + action + "'")
						os.Exit(0)
					}
				}
			}
		}
	}
	// Check connectors
	for t := range conf.Connectors {
		for e1 := range conf.Connectors[t].ExecGraph.Edges {
			for e2 := range conf.Connectors[t].ExecGraph.Edges[e1] {
				action := conf.Connectors[t].ExecGraph.Edges[e1][e2].Action.ActionName
				if shared.IsExternal(action) {
					if action != shared.INVP && action != shared.TERP && action != shared.INVR && action != shared.TERR {
						fmt.Println("EE:: Connector '" + conf.Connectors[t].Id + "' has an invalid action: '" + action)
						os.Exit(0)
					}
				} else {
					if shared.IsInternal(action) {
						st := reflect.TypeOf(conf.Connectors[t].TypeElem)
						_, ok := st.MethodByName(action)
						if !ok {
							fmt.Println("EE: Connector '" + conf.Connectors[t].Id + "' has an invalid action: '" + action + "'")
							os.Exit(0)
						}

					} else {
						fmt.Println("EE: Connector '" + conf.Connectors[t].Id + "' has an invalid action: '" + action + "'")
						os.Exit(0)
					}
				}
			}
		}
	}
}

func (engine Engine) ConfigureManagementChannels(conf configuration.Configuration) map[string]chan commands.LowLevelCommand {
	managementChannels := make(map[string]chan commands.LowLevelCommand)
	for i := range conf.Components {
		id := conf.Components[i].Id
		managementChannels[id] = make(chan commands.LowLevelCommand)
	}
	return managementChannels
}

func (engine Engine) CreateExecGraphs(conf *configuration.Configuration) {

	// Components
	for c := range conf.Components {
		graph := execgraph.NewGraph(conf.Components[c].FDRGraph.NumNodes)
		eActions := execgraph.Action{}
		var msg messages.SAMessage
		for e1 := range conf.Components[c].FDRGraph.Edges {
			for e2 := range conf.Components[c].FDRGraph.Edges[e1] {
				edgeTemp := conf.Components[c].FDRGraph.Edges[e1][e2]
				actionNameFDR := edgeTemp.Action
				actionNameExec := ""
				if strings.Contains(actionNameFDR, ".") {
					actionNameExec = actionNameFDR[:strings.Index(actionNameFDR, ".")]
				}
				if shared.IsExternal(actionNameExec) { // External action
					key := conf.Components[c].Id + "." + actionNameFDR
					channel := conf.StructuralChannels[key]
					params := execgraph.Action{}
					switch actionNameExec {
					case "InvR":
						invr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvR, Message: &msg, ActionChannel: &invr, ActionName: actionNameExec}
					case "TerR":
						terr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.TerR, Message: &msg, ActionChannel: &terr, ActionName: actionNameExec}
					case "InvP":
						invp := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvP, Message: &msg, ActionChannel: &invp, ActionName: actionNameExec}
					case "TerP":
						terp := channel
						params = execgraph.Action{ExternalAction: element.Element{}.TerP, Message: &msg, ActionChannel: &terp, ActionName: actionNameExec}
					}
					mapType := execgraph.Action{}
					mapType = params
					eActions = mapType
				}

				if shared.IsInternal(actionNameFDR) {
					msg := messages.SAMessage{}
					channel := make(chan messages.SAMessage)
					params := execgraph.Action{InternalAction: shared.Invoke, ActionName: actionNameFDR, Message: &msg, ActionChannel: &channel}
					mapType := params
					eActions = mapType
				}
				graph.AddEdge(edgeTemp.From, edgeTemp.To, eActions)
			}
		}
		tempComp := conf.Components[c]
		tempComp.SetExecGraph(graph)
		conf.Components[c] = tempComp
	}

	// Connectors
	for t := range conf.Connectors {
		graph := execgraph.NewGraph(conf.Connectors[t].FDRGraph.NumNodes)
		eActions := execgraph.Action{}
		var msg messages.SAMessage
		for e1 := range conf.Connectors[t].FDRGraph.Edges {
			for e2 := range conf.Connectors[t].FDRGraph.Edges[e1] {
				edgeTemp := conf.Connectors[t].FDRGraph.Edges[e1][e2]
				actionNameFDR := edgeTemp.Action
				actionNameExec := ""
				if strings.Contains(actionNameFDR, ".") {
					actionNameExec = actionNameFDR[:strings.Index(actionNameFDR, ".")]
				}
				if shared.IsExternal(actionNameExec) { // External action
					key := conf.Connectors[t].Id + "." + actionNameFDR
					channel := conf.StructuralChannels[key]
					params := execgraph.Action{}
					switch actionNameExec {
					case "InvR":
						invr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvR, Message: &msg, ActionChannel: &invr, ActionName: actionNameExec}
					case "TerR":
						terr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.TerR, Message: &msg, ActionChannel: &terr, ActionName: actionNameExec}
					case "InvP":
						invp := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvP, Message: &msg, ActionChannel: &invp, ActionName: actionNameExec}
					case "TerP":
						terp := channel
						params = execgraph.Action{ExternalAction: element.Element{}.TerP, Message: &msg, ActionChannel: &terp, ActionName: actionNameExec}
					}
					mapType := execgraph.Action{}
					mapType = params
					eActions = mapType
				}

				if shared.IsInternal(actionNameFDR) {
					msg := messages.SAMessage{}
					params := execgraph.Action{InternalAction: shared.Invoke, ActionName: actionNameFDR, Message: &msg}
					mapType := params
					eActions = mapType
				}
				graph.AddEdge(edgeTemp.From, edgeTemp.To, eActions)
			}
		}
		tempComp := conf.Connectors[t]
		tempComp.SetExecGraph(graph)
		conf.Connectors[t] = tempComp
	}
}

func (Engine) ConfigureStructuralChannelsAndMaps(conf *configuration.Configuration) {
	structuralChannels := make(map[string]chan messages.SAMessage)

	// Configure structural channels
	for i := range conf.Attachments {
		c1Id := conf.Attachments[i].C1.Id
		c2Id := conf.Attachments[i].C2.Id
		tId := conf.Attachments[i].T.Id

		// c1 -> t
		key01 := c1Id + "." + shared.INVR + "." + tId
		key02 := tId + "." + shared.INVP + "." + c1Id
		key03 := tId + "." + shared.TERP + "." + c1Id
		key04 := c1Id + "." + shared.TERR + "." + tId
		structuralChannels[key01] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key02] = structuralChannels[key01]
		structuralChannels[key03] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key04] = structuralChannels[key03]

		// t -> c2
		key01 = tId + "." + shared.INVR + "." + c2Id
		key02 = c2Id + "." + shared.INVP + "." + tId
		key03 = c2Id + "." + shared.TERP + "." + tId
		key04 = tId + "." + shared.TERR + "." + c2Id
		structuralChannels[key01] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key02] = structuralChannels[key01]
		structuralChannels[key03] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key04] = structuralChannels[key03]
	}
	conf.StructuralChannels = structuralChannels

	// Configure maps
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
	conf.Maps = elemMaps
}
