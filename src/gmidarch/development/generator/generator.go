package generator

import (
	"gmidarch/development/artefacts"
	"errors"
	"gmidarch/shared/parameters"
	"strings"
	"gmidarch/shared/shared"
)

type Generator struct{}

func (g Generator) GenerateCSP(madlGo artefacts.MADLGo, maps map[string]string) (artefacts.CSP, error) {
	r1 := artefacts.CSP{}
	r2 := *new(error)

	// Generate CSP
	r1,err := generateCSP(madlGo,maps)
	if err != nil {
		r2 = errors.New("Generator:: " + err.Error())
		return r1,r2
	}

	return r1, r2
}

func generateCSP(madlGo artefacts.MADLGo, maps map[string]string) (artefacts.CSP,error){
	r1 := artefacts.CSP{}
	r2 := *new (error)

	// Generate CSP Mid
	r1, err := artefacts.CSP{}.Create(madlGo,maps)
	if err != nil {
		r2 = errors.New("EE:: " + err.Error())
		return r1,r2
	}

	return r1,r2
}

func (Generator) GenerateCSPFile(csp artefacts.CSP)(error){
	r1 := *new(error)

	// Generate File
	cspFile := artefacts.CSPFile{}

	// File Path
	cspFile.FilePath = parameters.DIR_CSP + "/" + csp.CompositionName

	// File Name
	cspFile.FileName = csp.CompositionName + parameters.CSP_EXTENSION

	// File content
	dataTypeExp := "datatype PROCNAMES = " + shared.StringComposition(csp.Datatype, "|", true)
	eChannelExp := "channel " + shared.StringComposition(csp.EChannels, ",", false) + " : PROCNAMES"
	iChannelExp := "channel " + shared.StringComposition(csp.IChannels, ",", false)
	processesExp := ""
	for i := range csp.CompProcesses {
		processesExp += csp.CompProcesses[i]+"\n"
	}
	for i := range csp.ConnProcesses {
		processesExp += csp.ConnProcesses[i]+ "\n"
	}

	compositionExp := csp.CompositionName + " = (" + strings.ToUpper(shared.StringComposition(csp.Composition.Components, "|||", true)+")")
	compositionExp += "[|{|" + shared.StringComposition(csp.Composition.SyncPorts, ",", false) + "|}|]"

	renamings := []string{}
	conns := []string{}
	for i := range csp.Composition.RenamingPorts {
		for j := range csp.Composition.RenamingPorts[i] {
			r := csp.Composition.RenamingPorts[i][j].OldName + " <- " + csp.Composition.RenamingPorts[i][j].NewName
			renamings = append(renamings, r)
		}
		conns = append(conns, strings.ToUpper(i)+"[["+shared.StringComposition(renamings, ",", false)+"]]")
	}
	compositionExp += "(" + shared.StringComposition(conns, "|||", true) + ")"
	propertyExp := shared.StringComposition(csp.Property, "\n", false)

	cspFile.Content = append(cspFile.Content, dataTypeExp+"\n")
	cspFile.Content = append(cspFile.Content, eChannelExp+"\n")
	cspFile.Content = append(cspFile.Content, iChannelExp+"\n")
	cspFile.Content = append(cspFile.Content, processesExp+"\n")
	cspFile.Content = append(cspFile.Content, compositionExp+"\n")
	cspFile.Content = append(cspFile.Content, propertyExp)

	// Save file
	r1 = cspFile.Save()
	if r1 != nil {
		r1 = errors.New("Generator:: "+r1.Error())
		return r1
	}
	return r1
}