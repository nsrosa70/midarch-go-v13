package artefacts

import (
	"shared/parameters"
	"strings"
	"shared/shared"
	"errors"
)

type CompositionProcess struct {
	Components    []string
	Connectors    [] string
	SyncPorts     [] string
	RenamingPorts map[string][]Renaming
}

type Renaming struct {
	OldName string
	NewName string
}

type CSP struct {
	Maps            map[string]string
	CompositionName string
	Datatype        [] string
	IChannels       []string
	EChannels       []string
	CompProcesses   map[string]string
	ConnProcesses   map[string]string
	Composition     CompositionProcess
	Property        []string
}

func (c CSP) GenerateFile() CSPFile {
	r1 := CSPFile{}

	// File Path
	r1.FilePath = parameters.DIR_CSP

	// File Name
	r1.FileName = c.CompositionName + parameters.CSP_EXTENSION

	// File content
	dataTypeExp := "datatype PROCNAMES = " + c.StringComposition(c.Datatype, "|", true)
	eChannelExp := "channel " + c.StringComposition(c.EChannels, ",", false)
	iChannelExp := "channel " + c.StringComposition(c.IChannels, ",", false) + " : PROCNAMES"
	processesExp := ""
	for i := range c.CompProcesses {
		processesExp += strings.Replace(c.CompProcesses[i], "B", strings.ToUpper(i), 99) + "\n"
	}
	for i := range c.ConnProcesses {
		processesExp += strings.Replace(c.ConnProcesses[i], "B", strings.ToUpper(i), 99) + "\n"
	}

	compositionExp := c.CompositionName + " = (" + strings.ToUpper(c.StringComposition(c.Composition.Components, "|||", true)+")")
	compositionExp += "[|{|" + c.StringComposition(c.Composition.SyncPorts, ",", false) + "|}|]"

	renamings := []string{}
	conns := []string{}
	for i := range c.Composition.RenamingPorts {
		for j := range c.Composition.RenamingPorts[i] {
			r := c.Composition.RenamingPorts[i][j].OldName + " <- " + c.Composition.RenamingPorts[i][j].NewName
			renamings = append(renamings, r)
		}
		conns = append(conns, strings.ToUpper(i)+"[["+c.StringComposition(renamings, ",", false)+"]]")
	}
	compositionExp += "(" + c.StringComposition(conns, "|||", true) + ")"
	propertyExp := c.StringComposition(c.Property, "\n", false)

	r1.Content = append(r1.Content, dataTypeExp+"\n")
	r1.Content = append(r1.Content, eChannelExp+"\n")
	r1.Content = append(r1.Content, iChannelExp+"\n")
	r1.Content = append(r1.Content, processesExp+"\n")
	r1.Content = append(r1.Content, compositionExp+"\n")
	r1.Content = append(r1.Content, propertyExp)

	return r1
}

func (c *CSP) Create(madlGo MADLGO) error {
	r1 := *new(error)

	// Configure Maps
	c.Maps = madlGo.Maps

	// File name
	c.CompositionName = madlGo.ConfigurationName

	// Data type
	dataTypes := []string{}
	for c := range madlGo.Components {
		dataTypes = append(dataTypes, madlGo.Components[c].ElemId)
	}
	for t := range madlGo.Connectors {
		dataTypes = append(dataTypes, madlGo.Connectors[t].ElemId)
	}
	c.Datatype = dataTypes

	// Internal Channels
	c.IChannels = c.IdentifyInternalChannels(madlGo)

	// External Channels
	c.EChannels = c.IdentifyExternalChannels(madlGo)

	// Processes
	compProcesses := map[string]string{}
	for c := range madlGo.Components {
		compProcesses[madlGo.Components[c].ElemId] = madlGo.Components[c].CSP
	}
	c.CompProcesses = compProcesses

	connProcesses := map[string]string{}
	for t := range madlGo.Connectors {
		connProcesses[madlGo.Connectors[t].ElemId] = madlGo.Connectors[t].CSP
	}
	c.ConnProcesses = connProcesses

	// Processes - Configure Process Behaviours
	c.ConfigureProcessBehaviours()

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
	for i := range c.EChannels {
		cannonicalName, r1 := c.ToCanonicalName(c.EChannels[i])
		if r1 != nil {
			r1 = errors.New("CSP:: ")
		}
		cannonicalNames[cannonicalName] = cannonicalName
	}
	for i := range cannonicalNames {
		compositionTemp.SyncPorts = append(compositionTemp.SyncPorts, cannonicalNames[i])
	}

	// Composition Process - Renaming ports
	eChannels := map[string][]string{}
	for i := range c.ConnProcesses {
		tokens := c.MyTokenize(c.ConnProcesses[i])
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
			renaming := Renaming{OldName: eChannels[i][j], NewName: c.RenameSyncPort(eChannels[i][j], i)}
			renamingPorts = append(renamingPorts, renaming)
		}
		compositionTemp.RenamingPorts[i] = renamingPorts
	}
	c.Composition = compositionTemp

	// Property
	c.Property = append(c.Property, strings.Replace(parameters.DEADLOCK_PROPERTY, parameters.CORINGA, madlGo.ConfigurationName, 99))

	return r1
}

func (c CSP) IdentifyInternalChannels(madl MADLGO) []string {
	r1 := []string{}
	r1Temp := map[string]string{}

	for i := range madl.Components {
		tokens := c.MyTokenize(madl.Components[i].CSP)
		for j := range tokens {
			if shared.IsInternal(tokens[j]) {
				iAction := strings.TrimSpace(tokens[j])
				r1Temp[iAction] = iAction
			}
		}
	}

	for i := range madl.Connectors {
		tokens := c.MyTokenize(madl.Connectors[i].CSP)
		for i := range tokens {
			if shared.IsInternal(tokens[i]) {
				iAction := strings.TrimSpace(tokens[i])
				r1Temp[iAction] = iAction
			}
		}
	}

	for i := range r1Temp {
		r1 = append(r1, i)
	}
	return r1
}

func (c CSP) IdentifyExternalChannels(madl MADLGO) []string {
	r1 := []string{}
	r1Temp := map[string]string{}

	for i := range madl.Components {
		tokens := c.MyTokenize(madl.Components[i].CSP)
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				iAction := strings.TrimSpace(tokens[j])
				r1Temp[iAction] = iAction
			}
		}
	}

	for i := range madl.Connectors {
		tokens := c.MyTokenize(madl.Connectors[i].CSP)
		for i := range tokens {
			if shared.IsExternal(tokens[i]) {
				iAction := strings.TrimSpace(tokens[i])
				r1Temp[iAction] = iAction
			}
		}
	}

	for i := range r1Temp {
		r1 = append(r1, i)
	}
	return r1
}

func (CSP) ToCanonicalName(name string) (string, error) {
	r1 := ""
	r2 := *new(error)

	if strings.Contains(name, parameters.INVR) {
		r1 = parameters.INVR
	}
	if strings.Contains(name, parameters.TERR) {
		r1 = parameters.TERR
	}
	if strings.Contains(name, parameters.INVP) {
		r1 = parameters.INVP
	}
	if strings.Contains(name, parameters.TERP) {
		r1 = parameters.TERP
	}

	if r1 == "" {
		r2 = errors.New("Channel '" + name + "' has NOT a cannonical name.")
	}

	return r1, r2
}

func (CSP) MyTokenize(s string) [] string {
	tokens := []string{}

	token := ""
	for i := 0; i < len(s); i++ {
		c := s[i : i+1]
		switch c {
		case "=":
			token = ""
		case "-":
			if strings.TrimSpace(token) != "" {
				tokens = append(tokens, token)
			}
			token = ""
		case " ":
			if strings.TrimSpace(token) != "" {
				tokens = append(tokens, token)
			}
			token = ""
		case "]":
			token = ""
		case ">":
			token = ""
		case "\n":
			token = ""
		case "[":
			if strings.TrimSpace(token) != "" {
				tokens = append(tokens, token)
			}
			token = ""
		case "(":
			if strings.TrimSpace(token) != "" {
				tokens = append(tokens, token)
			}
			token = ""
		case ")":
			if strings.TrimSpace(token) != "" {
				tokens = append(tokens, token)
			}
			token = ""
		default:
			token += c
		}
	}
	return tokens
}

func (CSP) StringComposition(e []string, sep string, hasSpace bool) string {
	r1 := ""

	for i := range e {
		if hasSpace {
			r1 += e[i] + " " + sep + " "
		} else {
			r1 += e[i] + sep
		}
	}
	if hasSpace {
		r1 = r1[:len(r1)-len(sep)-2]
	} else {
		r1 = r1[:len(r1)-len(sep)]
	}

	return r1
}

func (c *CSP) ConfigureProcessBehaviours() {

	// Components
	for i := range c.CompProcesses {
		configuredBehaviour := c.CompProcesses[i]
		tokens := strings.Split(configuredBehaviour, " ")
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				eX := tokens[j][strings.Index(tokens[j], ".")+1:]
				key := i + "." + eX
				partner := c.Maps[key]
				configuredBehaviour = strings.Replace(configuredBehaviour, eX, partner, 99)
			}
		}
		c.CompProcesses[i] = configuredBehaviour
	}

	// Connectors
	for i := range c.ConnProcesses {
		configuredBehaviour := c.ConnProcesses[i]
		tokens := strings.Split(configuredBehaviour, " ")
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				eX := tokens[j][strings.Index(tokens[j], ".")+1:]
				key := i + "." + eX
				partner := c.Maps[key]
				configuredBehaviour = strings.Replace(configuredBehaviour, eX, partner, 99)
			}
		}
		c.ConnProcesses[i] = configuredBehaviour
	}
}

func (CSP) RenameSyncPort(action string, processId string) string {
	r1 := ""

	action = action [0:strings.Index(action, ".")]
	switch action {
	case parameters.INVP:
		r1 = parameters.INVR + "." + processId
	case parameters.TERP:
		r1 = parameters.INVR + "." + processId
	case parameters.INVR:
		r1 = parameters.INVR + "." + processId
	case parameters.TERR:
		r1 = parameters.INVR + "." + processId
	}
	return r1
}
