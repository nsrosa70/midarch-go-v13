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
)

type FDR struct{}

func (FDR) CheckBehaviour(conf configuration.Configuration, elemMaps map[string]string) bool {

	csp := createCSP(conf, elemMaps)
	saveCSP(conf, csp)
	r := invokeFDR(conf)

	return r
}

func saveCSP(conf configuration.Configuration, csp string) {

	// create file if not exists && truncate otherwise
	fileName := conf.Id + ".csp"
	file, err := os.Create(parameters.DIR_CSP + "/" + fileName)
	if err != nil {
		myError := errors.MyError{Source: "FDR", Message: "CSP File not created"}
		myError.ERROR()
	}
	defer file.Close()

	// save data
	_, err = file.WriteString(csp)
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
	cmdExp := commands.FDR_COMMAND
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
	generalBehaviour := createGeneralBehaviourExp(conf,externalChannels,componentProcesses)

	// assertion
	assertion := "assert P1 :[deadlock free]"
	csp := dataTypeExp + "\n" + internalChannelsExp + "\n" + externalChannelsExp + "\n" + processesExp + "\n" + generalBehaviour + "\n" + assertion

	return csp
}

func adjustPartnersComponent(id string, behaviour string) string {

	for i := 0; i < 99; i++ {
		behaviour = strings.Replace(behaviour, "e"+strconv.Itoa(i), id, 99)
	}
	return behaviour
}

func adjustPartnersConnectors(id string, behaviour string, elemMaps map[string]string) string {
	numPartners := strings.Count(behaviour,".e")
	for i := 0; i < numPartners; i++ {
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
	actions := shared.ToActions(behaviour)
	renamingExp := strings.ToUpper(id) + "[["
	for i := range actions {
		action := strings.TrimSpace(actions[i])
		if !strings.Contains(action, "I_") {
			action := action[0:strings.Index(action, ".")]
			switch action {
			case "InvP":
				renamingExp += "InvP <- InvR" + ","
			case "TerP":
				renamingExp += "TerP <- TerR" + ","
			case "InvR":
				renamingExp += "InvR <- InvP" + ","
			case "TerR":
				renamingExp += "TerR <- TerP" + ","
			}
		}
	}

	renamingExp = renamingExp[0:strings.LastIndex(renamingExp, ",")] + "]]"
	return renamingExp
}

func createDataTypeExp(conf configuration.Configuration) string {

	// datatypes
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
		behaviour := conf.Components[i].BehaviourExp
		actions := shared.ToActions(behaviour)
		for j := range actions {
			if strings.Contains(actions[j], "I_") {
				internalChannels[actions[j]] = strings.TrimSpace(actions[j])
			}
		}
	}

	for i := range conf.Connectors {
		behaviour := conf.Connectors[i].BehaviourExp
		actions := shared.ToActions(behaviour)
		for j := range actions {
			if strings.Contains(actions[j], "I_") {
				internalChannels[actions[j]] = strings.TrimSpace(actions[j])
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
		behaviour := conf.Components[i].BehaviourExp
		actions := shared.ToActions(behaviour)
		for j := range actions {
			if !strings.Contains(actions[j], "I_") {
				tempChannel := strings.TrimSpace(actions[j])
				tempChannel = tempChannel[0:strings.Index(tempChannel, ".")]
				externalChannels[tempChannel] = tempChannel
			}
		}
	}
	for i := range conf.Connectors {
		behaviour := conf.Connectors[i].BehaviourExp
		actions := shared.ToActions(behaviour)
		for j := range actions {
			if !shared.IsInternal(actions[j]) {
				tempChannel := strings.TrimSpace(actions[j])
				tempChannel = tempChannel[0:strings.Index(tempChannel, ".")]
				externalChannels[tempChannel] = tempChannel
			}
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
