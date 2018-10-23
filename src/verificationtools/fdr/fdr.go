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
	"framework/libraries"
	"reflect"
)

type FDR struct{}

func (FDR) CheckBehaviour(conf configuration.Configuration) bool {

	conf.CSP = createCSP(conf)
	saveCSP(conf)
	r := invokeFDR(conf)

	return r
}

func (FDR) LoadFDRGraph(conf *configuration.Configuration) {
	graph := fdrgraph.NewGraph(100)

	dotFileName := strings.Replace(conf.ADLFileName, ".conf", ".dot", 1)
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
	conf.FDRGraph = *graph
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

func createCSP(conf configuration.Configuration) string {

	// Configure behaviour expressions of components and conncetores
	configureBehaviours(&conf)

	// general behaviour
	dataTypeExp := createDataTypeExp(conf)
	internalChannelsExp, _ := createInternalChannelExp(conf)
	externalChannelsExp, externalChannels := createExternalChannelExp(conf)
	processesExp, componentProcesses, _ := createProcessExp(conf)
	generalBehaviour := createGeneralBehaviourExp(&conf, externalChannels, componentProcesses)

	// assertion
	assertion := "assert P1 :[deadlock free]"
	csp := dataTypeExp + "\n" + internalChannelsExp + "\n" + externalChannelsExp + "\n" + processesExp + "\n" + generalBehaviour + "\n" + assertion

	return csp
}

func configureBehaviours(conf *configuration.Configuration) {

	for i := range conf.Components {
		standardBehaviour := libraries.Repository[reflect.TypeOf(conf.Components[i].TypeElem).String()].CSP
		tokens := strings.Split(standardBehaviour, " ")
		for j := range tokens {
			if shared.IsAction(tokens[j]) && !shared.IsInternal(tokens[j]) {
				eX := tokens[j][strings.Index(tokens[j], ".")+1:]
				standardBehaviour = strings.Replace(standardBehaviour, eX, conf.Components[i].Id, 99)
			}
		}
		configuredBehaviour := strings.Replace(standardBehaviour, "B", strings.ToUpper(conf.Components[i].Id), 99)
		configuredBehaviour = shared.RenameInternalChannels(configuredBehaviour, conf.Components[i].Id)
		conf.Components[i] = element.Element{Id: conf.Components[i].Id, TypeElem: conf.Components[i].TypeElem, CSP: configuredBehaviour}
	}

	for i := range conf.Connectors {
		standardBehaviour := libraries.Repository[reflect.TypeOf(conf.Connectors[i].TypeElem).String()].CSP
		tokens := strings.Split(standardBehaviour, " ")
		for j := range tokens {
			if shared.IsAction(tokens[j]) && !shared.IsInternal(tokens[j]) {
				partner := tokens[j][strings.Index(tokens[j], ".")+1:]
				key := conf.Connectors[i].Id + "." + partner
				standardBehaviour = strings.Replace(standardBehaviour, partner, conf.Maps[key], 99)
			}
		}
		configuredBehaviour := strings.Replace(standardBehaviour, "B", strings.ToUpper(conf.Connectors[i].Id), 99)
		configuredBehaviour = shared.RenameInternalChannels(configuredBehaviour, conf.Connectors[i].Id)
		conf.Connectors[i] = element.Element{Id: conf.Connectors[i].Id, TypeElem: conf.Connectors[i].TypeElem, CSP: configuredBehaviour}
	}
}

func renamingSyncPorts(conf *configuration.Configuration, elem element.Element) string {
	id := elem.Id
	standardBehaviour := libraries.Repository[reflect.TypeOf(conf.Connectors[id].TypeElem).String()].CSP
	tokens := strings.Split(standardBehaviour, " ")
	renamingExp := strings.ToUpper(id) + "[["

	for i := range tokens {
		token := strings.TrimSpace(tokens[i])
		if !shared.IsInternal(token) && shared.IsAction(token) {
			action := token[0:strings.Index(token, ".")]
			eX := token[strings.Index(token, ".")+1:]
			key := id+"."+eX
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
			condiguredBehaviour := strings.Replace(standardBehaviour, eX, conf.Maps[key], 99)
			conf.Connectors[id] = element.Element{Id: conf.Connectors[id].Id, TypeElem: conf.Connectors[id].TypeElem, CSP: condiguredBehaviour}
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
		tokens := strings.Split(conf.Components[i].CSP, " ")
		for i := range tokens {
			if shared.IsInternal(tokens[i]) {
				iAction := strings.TrimSpace(tokens[i])
				internalChannels[iAction] = iAction
			}
		}
	}

	for i := range conf.Connectors {
		tokens := strings.Split(conf.Connectors[i].CSP, " ")
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
		b := conf.Components[i].CSP
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
		b := conf.Connectors[i].CSP
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

func createProcessExp(conf configuration.Configuration) (string, map[string]string, map[string]string) {

	componentProcesses := make(map[string]string)
	processesExp := ""
	for i := range conf.Components {
		id := conf.Components[i].Id
		componentProcesses[strings.ToUpper(id)] = conf.Components[i].CSP
		processesExp += componentProcesses[strings.ToUpper(id)] + "\n"
	}
	connectorProcesses := make(map[string]string)
	for i := range conf.Connectors {
		id := conf.Connectors[i].Id
		connectorProcesses[strings.ToUpper(id)] = conf.Connectors[i].CSP
		processesExp += connectorProcesses[strings.ToUpper(id)] + "\n"
	}
	return processesExp, componentProcesses, connectorProcesses
}

func createGeneralBehaviourExp(conf *configuration.Configuration, externalChannels map[string]string, componentProcesses map[string]string) string {

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
		generalBehaviour += renamingSyncPorts(conf, conf.Connectors[i]) + "|||"
	}
	generalBehaviour = generalBehaviour[0:strings.LastIndex(generalBehaviour, "|||")] + ")"

	return generalBehaviour
}
