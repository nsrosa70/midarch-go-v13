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

func (FDR) LoadFDRGraphs(conf *configuration.Configuration) {

	// Load component
	dotDir := parameters.DIR_CSP+"/"+strings.Replace(conf.Id,".conf","",99)
	for c := range conf.Components {
		dotFileName := strings.ToUpper(conf.Components[c].Id) + ".dot"
		dotFileName = dotDir + "/" + dotFileName

		fileContent := loadDotFile(dotFileName)
		graph := createFDRGraph(fileContent)
		tempComp := conf.Components[c]
		tempComp.SetFDRGraph(*graph)
		conf.Components[c] = tempComp
	}

	// Load connectors
	for t := range conf.Connectors {
		dotFileName := strings.ToUpper(conf.Connectors[t].Id) + ".dot"
		dotFileName = dotDir + "/" + dotFileName

		fileContent := loadDotFile(dotFileName)
		graph := createFDRGraph(fileContent)
		tempConn := conf.Connectors[t]
		tempConn.SetFDRGraph(*graph)
		conf.Connectors[t] = tempConn
	}
}

func createFDRGraph(fileContent []string) *fdrgraph.Graph {
	graph := fdrgraph.NewGraph(100)

	for l := range fileContent {
		line := fileContent[l]
		if strings.Contains(line, "->") {
			from, _ := strconv.Atoi(strings.TrimSpace(line[:strings.Index(line, "->")]))
			to, _ := strconv.Atoi(strings.TrimSpace(line[strings.Index(line, "->")+2 : strings.Index(line, "[")]))
			action := line[strings.Index(line, "=")+2 : strings.LastIndex(line, "]")-2]
			graph.AddEdge(from, to, action)
		}
	}

	return graph
}

func loadDotFile(dotFileName string) []string {

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

	return fileContent
}

func saveCSP(conf configuration.Configuration) {

	// create diretcory if it does not exist
	confDir := parameters.DIR_CSP+"/"+strings.Replace(conf.Id,".conf","",99)
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		os.MkdirAll(confDir, os.ModePerm);
	}

	// create file if it does not exist && truncate otherwise
	fileName := conf.Id + ".csp"
	file, err := os.Create(confDir + "/" + fileName)
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
	dirFile := parameters.DIR_CSP+"/"+strings.Replace(conf.Id,".conf","",99)
	inputFile := dirFile + "/" + fileName

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

	// Configure components
	for c := range conf.Components {
		standardBehaviour := libraries.Repository[reflect.TypeOf(conf.Components[c].TypeElem).String()].CSP
		tokens := strings.Split(standardBehaviour, " ")
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				eX := tokens[j][strings.Index(tokens[j], ".")+1:]
				key := conf.Components[c].Id + "." + eX
				partner := conf.Maps[key]
				standardBehaviour = strings.Replace(standardBehaviour, eX, partner, 99)
			}
		}
		configuredBehaviour := strings.Replace(standardBehaviour, "B", strings.ToUpper(conf.Components[c].Id), 99)
		conf.Components[c] = element.Element{Id: conf.Components[c].Id, TypeElem: conf.Components[c].TypeElem, CSP: configuredBehaviour}
	}

	// Configure connectors
	for t := range conf.Connectors {
		standardBehaviour := libraries.Repository[reflect.TypeOf(conf.Connectors[t].TypeElem).String()].CSP
		tokens := strings.Split(standardBehaviour, " ")
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				eX := tokens[j][strings.Index(tokens[j], ".")+1:]
				key := conf.Connectors[t].Id + "." + eX
				partner := conf.Maps[key]
				standardBehaviour = strings.Replace(standardBehaviour, eX, partner, 99)
			}
		}
		configuredBehaviour := strings.Replace(standardBehaviour, "B", strings.ToUpper(conf.Connectors[t].Id), 99)
		conf.Connectors[t] = element.Element{Id: conf.Connectors[t].Id, TypeElem: conf.Connectors[t].TypeElem, CSP: configuredBehaviour}
	}
}

func renameSyncPorts(conf *configuration.Configuration, elem element.Element) string {
	id := elem.Id
	tokens := strings.Split(elem.CSP, " ")
	renamingExp := strings.ToUpper(id) + "[["

	for i := range tokens {
		token := strings.TrimSpace(tokens[i])
		if shared.IsExternal(token) {
			action := token[0:strings.Index(token, ".")]
			switch action {
			case shared.INVP:
				renamingExp += token + " <- " + shared.INVR + "." + id + ","
			case shared.TERP:
				renamingExp += token + " <- " + shared.TERR + "." + id + ","
			case shared.INVR:
				renamingExp += token + " <- " + shared.INVP + "." + id + ","
			case shared.TERR:
				renamingExp += token + " <- " + shared.TERP + "." + id + ","
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
		generalBehaviour += renameSyncPorts(conf, conf.Connectors[i]) + "|||"
	}
	generalBehaviour = generalBehaviour[0:strings.LastIndex(generalBehaviour, "|||")] + ")"

	return generalBehaviour
}
