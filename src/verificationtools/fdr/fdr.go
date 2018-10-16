package fdr

import (
	"strings"
	"strconv"
	"os"
	"framework/element"
	"os/exec"
	"framework/configuration/configuration"
	"shared/shared"
	"shared/parameters"
	"shared/errors"
	"framework/configuration/commands"
	"graph/fdrgraph"
	"log"
	"bufio"
)

type FDR struct{}

func (FDR) CheckBehaviour(conf configuration.Configuration, elemMaps map[string]string) bool {

	conf.CSP = createCSP(conf, elemMaps)
	saveCSP(conf)
	r := invokeFDR(conf)

	return r
}

func (FDR) LoadFDRGraph(confFile string) fdrgraph.Graph {
	graph := fdrgraph.NewGraph(100)

	dotFileName := strings.Replace(confFile, ".conf", ".dot", 1)
	dotFileName = parameters.DIR_CSP + "/" + dotFileName

	// read file
	fileContent := []string{}
	file, err := os.Open(dotFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	// Generate Configuration
	for l := range fileContent {
		line := fileContent[l]
		if strings.Contains(line, "->") {
			from, _ := strconv.Atoi(strings.TrimSpace(line[:strings.Index(line, "->")]))
			to, _ := strconv.Atoi(strings.TrimSpace(line[strings.Index(line, "->")+2 : strings.Index(line, "[")]))
			action := line[strings.Index(line, "=")+2 : strings.LastIndex(line, "]")-2]
			graph.AddEdge(from, to, action)
		}
	}
	return *graph
}

func saveCSP(conf configuration.Configuration) {

	// create file if not exists && truncate otherwise
	fileName := conf.Id + ".csp"
	file, err := os.Create(parameters.DIR_CSP + "/" + fileName)
	if err != nil {
		myError := errors.MyError{Source: "FDR", Message: "CSP File not created"}
		myError.ERROR()
	}
	defer file.Close()

	// save data
	_, err = file.WriteString(conf.CSP)
	if err != nil {
		myError := errors.MyError{Source: "FDR", Message: "CSP Specification not saved"}
		myError.ERROR()
	}
	err = file.Sync()
	if err != nil {
		myError := errors.MyError{Source: "FDR", Message: "CSP File not Synced"}
		myError.ERROR()
	}
	defer file.Close()
}

func invokeFDR(conf configuration.Configuration) bool {
	cmdExp := parameters.DIR_FDR + "/" + commands.FDR_COMMAND
	fileName := conf.Id + ".csp"
	inputFile := parameters.DIR_CSP + "/" + fileName

	out, err := exec.Command(cmdExp, inputFile).Output()
	if err != nil {
		myError := errors.MyError{Source: "FDR", Message: "Problem in invoking FDR (e.g.,syntax error)"}
		myError.ERROR()
	}
	s := string(out[:])

	if strings.Contains(s, "Passed") {
		return true
	} else {
		myError := errors.MyError{Source: "FDR", Message: "Deadlock detected"}
		myError.ERROR()
		return false
	}
}

func createCSP(conf configuration.Configuration, elemMaps map[string]string) string {

	// general behaviour
	dataTypeExp := createDataTypeExp(conf)
	internalChannelsExp, _ := createInternalChannelExp(conf)
	externalChannelsExp, externalChannels := createExternalChannelExp(conf)
	processesExp, componentProcesses, _ := createProcessExp(conf, elemMaps)
	generalBehaviour := createGeneralBehaviourExp(conf, externalChannels, componentProcesses)

	// assertion
	assertion := "assert P1 :[deadlock free]"
	csp := dataTypeExp + "\n" + internalChannelsExp + "\n" + externalChannelsExp + "\n" + processesExp + "\n" + generalBehaviour + "\n" + assertion

	return csp
}

func adjustPartnersComponent(id string, behaviour string) string {
	numPartners := strings.Count(behaviour, ".e")

	for i := 1; i < numPartners+1; i++ {
		behaviour = strings.Replace(behaviour, "e"+strconv.Itoa(i), id, numPartners+1)
	}
	return behaviour
}

func adjustPartnersConnectors(id string, behaviour string, elemMaps map[string]string) string {
	numPartners := strings.Count(behaviour, ".e")

	for i := 1; i < numPartners+1; i++ {
		key := id + "." + "e" + strconv.Itoa(i)
		value, ok := elemMaps[key]
		if ok {
			behaviour = strings.Replace(behaviour, "e"+strconv.Itoa(i), value, numPartners)
		}
	}
	return behaviour
}

func renamingPorts(elem element.Element) string {
	id := elem.Id
	behaviour := elem.BehaviourExp
	tokens := strings.Split(behaviour," ")
	renamingExp := strings.ToUpper(id) + "[["

	for i := range tokens {
		token := strings.TrimSpace(tokens[i])
		if !shared.IsInternal(token) && shared.IsAction(token){
			action := token[0:strings.Index(token, ".")]
			switch action {
			case shared.INVP:
				renamingExp += shared.INVP + " <- " + shared.INVR + ","
			case shared.TERP:
				renamingExp += shared.TERP + " <- " + shared.TERR + ","
			case shared.INVR:
				renamingExp += shared.INVR + " <- " + shared.INVP + ","
			case shared.TERR:
				renamingExp += shared.TERR + " <- " + shared.TERP + ","
			}
		}
	}
	renamingExp = renamingExp[0:strings.LastIndex(renamingExp, ",")] + "]]"
	return renamingExp
}

func createDataTypeExp(conf configuration.Configuration) string {
	dataTypes := make(map[string]string)

	for i := range conf.Components {
		dataTypes [conf.Components[i].Id] = conf.Components[i].Id
	}
	for i := range conf.Connectors {
		dataTypes [conf.Connectors[i].Id] = conf.Connectors[i].Id
	}

	dataTypeExp := "datatype PROCNAMES = "
	for i := range dataTypes {
		dataTypeExp += dataTypes[i] + " | "
	}
	dataTypeExp = dataTypeExp[0:strings.LastIndex(dataTypeExp, "|")]

	return dataTypeExp
}

func createInternalChannelExp(conf configuration.Configuration) (string, map[string]string) {
	internalChannels := make(map[string]string)

	for i := range conf.Components {
		tokens := strings.Split(conf.Components[i].BehaviourExp, " ")
		for i := range tokens {
			if shared.IsInternal(tokens[i]) {
				iAction := strings.TrimSpace(tokens[i])
				internalChannels[iAction] = iAction
			}
		}
	}

	for i := range conf.Connectors {
		tokens := strings.Split(conf.Connectors[i].BehaviourExp, " ")
		for i := range tokens {
			if shared.IsInternal(tokens[i]) {
				iAction := strings.TrimSpace(tokens[i])
				internalChannels[iAction] = iAction
			}
		}
	}
	internalChannelsExp := "channel "
	for i := range internalChannels {
		internalChannelsExp += internalChannels[i] + ","
	}
	internalChannelsExp = internalChannelsExp[0:strings.LastIndex(internalChannelsExp, ",")]

	return internalChannelsExp, internalChannels
}

func createExternalChannelExp(conf configuration.Configuration) (string, map[string]string) {
	externalChannels := make(map[string]string)

	for i := range conf.Components {
		b := conf.Components[i].BehaviourExp
		if strings.Contains(b, shared.INVR) {
			externalChannels[shared.INVR] = shared.INVR
		}
		if strings.Contains(b, shared.TERR) {
			externalChannels[shared.TERR] = shared.TERR
		}
		if strings.Contains(b, shared.INVP) {
			externalChannels[shared.INVP] = shared.INVP
		}
		if strings.Contains(b, shared.TERP) {
			externalChannels[shared.TERP] = shared.TERP
		}
	}
	for i := range conf.Connectors {
		b := conf.Connectors[i].BehaviourExp
		if strings.Contains(b, shared.INVR) {
			externalChannels[shared.INVR] = shared.INVR
		}
		if strings.Contains(b, shared.TERR) {
			externalChannels[shared.TERR] = shared.TERR
		}
		if strings.Contains(b, shared.INVP) {
			externalChannels[shared.INVP] = shared.INVP
		}
		if strings.Contains(b, shared.TERP) {
			externalChannels[shared.TERP] = shared.TERP
		}
	}

	externalChannelsExp := "channel "
	for i := range externalChannels {
		externalChannelsExp += externalChannels[i] + ","
	}
	externalChannelsExp = externalChannelsExp[0:strings.LastIndex(externalChannelsExp, ",")] + " : PROCNAMES"

	return externalChannelsExp, externalChannels
}

func createProcessExp(conf configuration.Configuration, elemMaps map[string]string) (string, map[string]string, map[string]string) {

	componentProcesses := make(map[string]string)
	processesExp := ""
	for i := range conf.Components {
		id := conf.Components[i].Id
		behaviour := strings.Replace(conf.Components[i].BehaviourExp, "B", strings.ToUpper(id), 99)
		behaviour = adjustPartnersComponent(id, behaviour)
		componentProcesses[strings.ToUpper(id)] = behaviour
		processesExp += componentProcesses[strings.ToUpper(id)] + "\n"
	}
	connectorProcesses := make(map[string]string)
	for i := range conf.Connectors {
		id := conf.Connectors[i].Id
		behaviour := strings.Replace(conf.Connectors[i].BehaviourExp, "B", strings.ToUpper(id), 99)
		behaviour = adjustPartnersConnectors(id, behaviour, elemMaps)
		connectorProcesses[strings.ToUpper(id)] = behaviour
		processesExp += connectorProcesses[strings.ToUpper(id)] + "\n"
	}
	return processesExp, componentProcesses, connectorProcesses
}

func createGeneralBehaviourExp(conf configuration.Configuration, externalChannels map[string]string, componentProcesses map[string]string) string {

	// (C1 ||| C2 ||| C3)
	generalBehaviour := "P1 = ("
	for i := range componentProcesses {
		generalBehaviour += i + "|||"
	}
	generalBehaviour = generalBehaviour[0:strings.LastIndex(generalBehaviour, "|||")] + ")"

	generalBehaviour += "\n" + "[|{|"
	for i := range externalChannels {
		generalBehaviour += externalChannels[i] + ","
	}
	generalBehaviour = generalBehaviour[0:strings.LastIndex(generalBehaviour, ",")] + "|}|]" + "\n"

	generalBehaviour += "("
	for i := range conf.Connectors {
		generalBehaviour += renamingPorts(conf.Connectors[i]) + "|||"
	}
	generalBehaviour = generalBehaviour[0:strings.LastIndex(generalBehaviour, "|||")] + ")"

	return generalBehaviour
}