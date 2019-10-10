package csp

import (
	"errors"
	"strings"
	"newsolution/shared/shared"
	"newsolution/shared/parameters"
)

ty
import (
	"errors"
	"strings"
	"newsolution/shared/shared"
	"newsolution/shared/parameters"
)

func (CSP) Create(madlGo madl.MADLGo, maps map[string]string, kindOfMADL int, midAdaptability []string) (CSP, error) {
	r1 := CSP{}
	r2 := *new(error)

	// Solve RUNTIME behaviours
	r2 = r1.ConfigureProcessBehaviours(madlGo, maps, kindOfMADL, midAdaptability)
	if r2 != nil {
		r2 = errors.New("CSP" + r2.Error())
		return r1, r2
	}

	// File name
	r1.CompositionName = madlGo.ConfigurationName

	// Data type
	dataTypes := []string{}
	for c := range madlGo.Components {
		dataTypes = append(dataTypes, madlGo.Components[c].ElemId)
	}
	for t := range madlGo.Connectors {
		dataTypes = append(dataTypes, madlGo.Connectors[t].ElemId)
	}
	r1.Datatype = dataTypes

	// Internal Channels
	r1.IChannels = identifyInternalChannels(madlGo)

	// External Channels
	r1.EChannels = identifyExternalChannels(madlGo)

	// Processes - Components
	compProcesses := map[string]string{}
	for i := range madlGo.Components {
		comp := madlGo.Components[i]
		compId := strings.ToUpper(comp.ElemId)
		subprocesses := strings.Split(comp.CSP, "\n")
		if len(subprocesses) > 1 {
			renamedSubprocesses := renameSubprocesses(comp.CSP)
			compProcesses[compId] = strings.Replace(renamedSubprocesses, "B", compId, 99)
		} else {
			compProcesses[compId] = strings.Replace(comp.CSP, "B", compId, 99)
		}
	}
	r1.CompProcesses = compProcesses

	// Processes - Connectors
	connProcesses := map[string]string{}
	for t := range madlGo.Connectors {
		conn := madlGo.Connectors[t]
		connId := strings.ToUpper(madlGo.Connectors[t].ElemId)
		connProcesses[connId] = strings.Replace(conn.CSP, "B", connId, 99)
	}
	r1.ConnProcesses = connProcesses

	// Composition process - Components/Connectors
	compositionTemp := CompositionProcess{}
	for i := range madlGo.Components {
		compositionTemp.Components = append(compositionTemp.Components, madlGo.Components[i].ElemId)
	}
	for i := range madlGo.Connectors {
		compositionTemp.Connectors = append(compositionTemp.Connectors, madlGo.Connectors[i].ElemId)
	}

	// Composition Process - Sync ports
	cannonicalNames := map[string]string{}
	for i := range r1.EChannels {
		cannonicalName, r2 := toCanonicalName(r1.EChannels[i])
		if r2 != nil {
			r2 = errors.New("CSP:: " + r2.Error())
			return r1, r2
		}
		cannonicalNames[cannonicalName] = cannonicalName
	}
	for i := range cannonicalNames {
		compositionTemp.SyncPorts = append(compositionTemp.SyncPorts, cannonicalNames[i])
	}

	// Composition Process - Renaming ports
	eChannels := map[string][]string{}
	for i := range r1.ConnProcesses {
		tokens := shared.MyTokenize(r1.ConnProcesses[i])
		actions := []string{}
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				actions = append(actions, tokens[j])
			}
			eChannels [i] = actions
		}
	}

	compositionTemp.RenamingPorts = map[string][]Renaming{}
	for i := range eChannels {
		renamingPorts := []Renaming{}
		for j := range eChannels[i] {
			renaming := Renaming{OldName: eChannels[i][j], NewName: r1.RenameSyncPort(eChannels[i][j], i)}
			renamingPorts = append(renamingPorts, renaming)
		}
		compositionTemp.RenamingPorts[i] = renamingPorts
	}
	r1.Composition = compositionTemp

	// Property
	r1.Property = append(r1.Property, strings.Replace(parameters.DEADLOCK_PROPERTY, parameters.CORINGA, madlGo.ConfigurationName, 99))

	return r1, r2
}
