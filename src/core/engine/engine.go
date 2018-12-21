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
	"core/versioninginjector"
	"framework/components"
	"framework/libraries"
	"framework/connectors"
)

type Engine struct{}

func (engine Engine) Deploy(adlApp string) {

	// Initialize environment
	fmt.Print("Engine:: Initialization...")
	engine.Initialize()
	fmt.Println(" (ok)")

	// Prepare App configuration to be executed
	fmt.Print("Engine:: Prepapre App configuration...")
	appConf := engine.PrepareConfiguration(adlApp,true)
	fmt.Println(" (ok)")

	// Generate adl of the execution environment
	fmt.Print("Engine:: Genarate ADL of EE...")
	adlEE := parameters.PREFIX_ADL_EXECUTION_ENVIRONMENT+adlApp
	generateAdlEE(appConf,adlEE)
	fmt.Println(" (ok)")

	// prepare EE configuration to be executed
	fmt.Print("Engine:: Prepapre EE configuration...")
	eeConf := engine.PrepareConfiguration(adlEE,false)
	fmt.Println(" (ok)")

	// Configure 'Info' of components belonging to the Execution Environment
	fmt.Print("Engine:: Prepapre 'Info' of EE...")
	engine.ConfigureInfoEE(&eeConf, appConf)
	fmt.Println(" (ok)")

	// Configure 'Info' of components beloging to the Application
	fmt.Print("Engine:: Prepapre 'Info' of Application...")
	engine.ConfigureApp(&appConf)
	fmt.Println(" (ok)")

	// Deploy App into Units
	fmt.Print("Engine:: Deploy Units...")
	engine.ConfigureUnits(&eeConf, appConf)
	fmt.Println(" (ok)")

	// Start configuration
	fmt.Print("Engine:: Start Configuration...")
	engine.StartConfiguration(eeConf)
	fmt.Println(" (ok)")

	// start versioning (if the architecture is adaptable)
	fmt.Print("Engine:: Start Versioning...")
	go versioninginjector.InjectAdaptiveEvolution(parameters.PLUGIN_BASE_NAME)
	fmt.Println(" (ok)")
}

func (ee Engine) Initialize() {
	// Load execution parameters
	shared.LoadParameters(os.Args[1:])

	// Check the CSP behaviours
	libraries.CheckLibrary()

	// Show execution parameters
	shared.ShowExecutionParameters(false)
}

func (Engine) ConfigureInfoEE(eeConf *configuration.Configuration, appConf configuration.Configuration) {
	for c1 := range eeConf.Components {
		componentType := reflect.TypeOf(eeConf.Components[c1].TypeElem).String()
		switch  componentType {
		case reflect.TypeOf(components.ExecutionEnvironment{}).String(): // Execution environment
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
		case reflect.TypeOf(components.MAPEKPlanner{}).String(): // Planner
			elem := eeConf.Components[c1]
			elem.SetInfo(appConf.Components)
			eeConf.Components[c1] = elem
		default:
			elem := eeConf.Components[c1]
			elem.SetInfo(parameters.DEFAULT_INFO)
			eeConf.Components[c1] = elem
		}
	}
}

func (Engine) ConfigureApp(conf *configuration.Configuration) {
	for c := range conf.Components {
		elem := conf.Components[c]
		elem.SetInfo(parameters.DEFAULT_INFO)
		conf.Components[c] = elem
	}
}

func (Engine) ConfigureUnits(eeConf *configuration.Configuration, appConf configuration.Configuration) {
	availableUnits := []string{}

	// Identify units
	for u := range eeConf.Components {
		if reflect.TypeOf(eeConf.Components[u].TypeElem).String() == reflect.TypeOf(components.ExecutionUnit{}).String() {
			availableUnits = append(availableUnits, eeConf.Components[u].Id)
		}
	}

	// Check if the numbers of units is enough to accomodate the application
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

func (Engine) StartConfiguration(conf configuration.Configuration) {
	for c := range conf.Components {
		msg := messages.SAMessage{}
		for e1 := range conf.Components[c].ExecGraph.Edges {
			for e2 := range conf.Components[c].ExecGraph.Edges[e1] {
				conf.Components[c].ExecGraph.Edges[e1][e2].Action.Message = &msg
			}
		}
		go startElement(conf.Components[c])
	}

	for t := range conf.Connectors {
		msg := messages.SAMessage{}
		for e1 := range conf.Connectors[t].ExecGraph.Edges {
			for e2 := range conf.Connectors[t].ExecGraph.Edges[e1] {
				conf.Connectors[t].ExecGraph.Edges[e1][e2].Action.Message = &msg
			}
		}
		go startElement(conf.Connectors[t])
	}
}

func startElement(elem element.Element) {
	for {
		shared.Invoke(elem, "Loop", elem, elem.ExecGraph)
	}
}

func (engine Engine) PrepareConfiguration(adlApp string, checkConfiguration bool) configuration.Configuration {

	// Generate Go configuration
	conf := configuration.MapADLIntoGo(adlApp)

	// Configure structural channels and maps of components/connectors
	engine.ConfigureStructuralChannelsAndMaps(&conf)

	// Update the behaviour of some connectors according to the configuration, e.g., OneToN
	updateBehaviourOfNConnectors(conf)

	// Check behaviour using FDR
	fdrChecker := new(fdr.FDR)
	conf.CSP = fdrChecker.CreateCSP(conf)
	fdrChecker.SaveCSP(conf)
	if checkConfiguration {
		ok := fdrChecker.InvokeFDR(conf)
		if !ok {
			myError := errors.MyError{Source: "Execution Engine", Message: "Configuration has a problem detected by FDR4"}
			myError.ERROR()
		}
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

func CheckActionsAndImplementations(conf configuration.Configuration) {

	// Check components
	for c := range conf.Components {
		for e1 := range conf.Components[c].ExecGraph.Edges {
			for e2 := range conf.Components[c].ExecGraph.Edges[e1] {
				action := conf.Components[c].ExecGraph.Edges[e1][e2].Action.ActionName
				if shared.IsExternal(action) {
					if action != parameters.INVP && action != parameters.TERP && action != parameters.INVR && action != parameters.TERR {
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
					if action != parameters.INVP && action != parameters.TERP && action != parameters.INVR && action != parameters.TERR {
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
					case parameters.INVR:
						invr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvR, Message: &msg, ActionChannel: &invr, ActionName: actionNameExec}
					case parameters.TERR:
						terr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.TerR, Message: &msg, ActionChannel: &terr, ActionName: actionNameExec}
					case parameters.INVP:
						invp := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvP, Message: &msg, ActionChannel: &invp, ActionName: actionNameExec}
					case parameters.TERP:
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
					case parameters.INVR:
						invr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvR, Message: &msg, ActionChannel: &invr, ActionName: actionNameExec}
					case parameters.TERR:
						terr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.TerR, Message: &msg, ActionChannel: &terr, ActionName: actionNameExec}
					case parameters.INVP:
						invp := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvP, Message: &msg, ActionChannel: &invp, ActionName: actionNameExec}
					case parameters.TERP:
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
		key01 := c1Id + "." + parameters.INVR + "." + tId
		key02 := tId + "." + parameters.INVP + "." + c1Id
		key03 := tId + "." + parameters.TERP + "." + c1Id
		key04 := c1Id + "." + parameters.TERR + "." + tId
		structuralChannels[key01] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key02] = structuralChannels[key01]
		structuralChannels[key03] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key04] = structuralChannels[key03]

		// t -> c2
		key01 = tId + "." + parameters.INVR + "." + c2Id
		key02 = c2Id + "." + parameters.INVP + "." + tId
		key03 = c2Id + "." + parameters.TERP + "." + tId
		key04 = tId + "." + parameters.TERR + "." + c2Id
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

func generateAdlEE(appConf configuration.Configuration, adlEE string) {

	// Generate units
	units := []string{}
	componentUnits := ""
	fullTypeName := reflect.TypeOf(components.ExecutionUnit{}).String()
	unitTypeName := fullTypeName[strings.LastIndex(fullTypeName, ".")+1:]
	for i := 0; i < len(appConf.Components)+len(appConf.Connectors); i++ {
		units = append(units, "unit"+strconv.Itoa(i+1))
	}
	for i := 0; i < len(units); i++ {
		componentUnits += "    "+units[i] + " : " + unitTypeName + " \n"
	}

	// Generate Attachments
	attUnits := ""
	for i := 0; i < len(units); i++ {
		attUnits += "   ee,t1," + units[i] + " \n"
	}

	// Assemble configuration
	header := "Configuration " + strings.Replace(adlEE, ".conf", "", 99) + " := "
	adaptability := "Adaptability \n None"
	components := "Components \n" +
		"    ee : ExecutionEnvironment \n" +
		"    evolutiveMonitor : MAPEKEvolutiveMonitor \n" +
		"    mapekMonitor : MAPEKMonitor \n" +
		"    analyser : MAPEKAnalyser \n" +
		"    planner : MAPEKPlanner \n" +
		"    executor : MAPEKExecutor \n" +
		componentUnits

	connectors := "Connectors \n" +
		"    t1 : OneToN \n" +
		"    t2 : OneWay \n" +
		"    t3 : OneWay \n" +
		"    t4 : OneWay \n" +
		"    t5 : OneWay \n" +
		"    t6 : OneWay"

	attachments := "Attachments \n" +
		attUnits +
		"   evolutiveMonitor,t2,mapekMonitor\n" +
		"   mapekMonitor,t3,analyser\n" +
		"   analyser,t4,planner\n" +
		"   planner,t5,executor\n" +
		"   executor,t6,ee"
	endConf := "EndConf"

	adl := header + "\n\n" + adaptability + "\n\n" + components + "\n\n" + connectors + "\n\n" + attachments + "\n\n" + endConf

	file, err := os.Create(parameters.DIR_CONF + "/" + adlEE)
	if err != nil {
		myError := errors.MyError{Source: "Engine", Message: "File '" + adlEE + "' NOT created"}
		myError.ERROR()
	}
	defer file.Close()

	// save data
	_, err = file.WriteString(adl)
	if err != nil {
		myError := errors.MyError{Source: "Engine", Message: "File '" + adlEE + "' NOT saved"}
		myError.ERROR()
	}
}

func updateBehaviourOfNConnectors(conf configuration.Configuration){

	// Find current behaviour in the Repository
	for i := range conf.Connectors{
		if reflect.TypeOf(conf.Connectors[i].TypeElem) == reflect.TypeOf(connectors.OneToN{}) {
			oldRecord := libraries.Repository[reflect.TypeOf(conf.Connectors[i].TypeElem).String()]
			newRecord := oldRecord
			n := countAttachments(conf,i)
			newBehaviour := defineNewBehaviour(n,connectors.OneToN{})
			newRecord.SetCSP(newBehaviour)
			libraries.Repository[reflect.TypeOf(conf.Connectors[i].TypeElem).String()] = newRecord
		}
	}
}

func countAttachments(conf configuration.Configuration, connectorId string) int {
	n := 0
	for i := range conf.Attachments{
		if conf.Attachments[i].T.Id == connectorId {
			n ++
		}
	}
	return n
}

func defineNewBehaviour(n int,elem interface{}) string{
	baseBehaviour := ""

	switch reflect.TypeOf(elem).String() {
	case reflect.TypeOf(connectors.OneToN{}).String():
		baseBehaviour = "B = InvP.e1"
		for i := 0; i < n; i++{
			baseBehaviour += " -> InvR.e"+strconv.Itoa(i+2)
		}
		baseBehaviour += " -> B"
	default:
		fmt.Println("Configuration:: Impossible to define the new behaviour of "+reflect.TypeOf(elem).String())
		os.Exit(0)
	}
	return baseBehaviour
}