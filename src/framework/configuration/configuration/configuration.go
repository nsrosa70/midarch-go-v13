package configuration

import (
	"framework/element"
	"framework/configuration/attachments"
	"graph/fdrgraph"
	"graph/execgraph"
	"os"
	"strings"
	"fmt"
	"framework/libraries"
	"shared/shared"
	"shared/parameters"
	"log"
	"bufio"
	"framework/components/srh"
	"framework/components/crh"
	"framework/message"
)

type Configuration struct {
	ADLFileName        string
	Id                 string
	Components         map[string]element.Element
	Connectors         map[string]element.Element
	Attachments        [] attachments.Attachment
	CSP                string
	FDRGraph           fdrgraph.Graph
	StateMachine       execgraph.Graph
	StructuralChannels map[string]chan message.Message
	Maps               map[string]string
}

func (conf *Configuration) AddComp(comp element.Element) {

	if conf.Components == nil {
		conf.Components = make(map[string]element.Element)
		conf.Components[comp.Id] = comp
	} else {
		conf.Components[comp.Id] = comp
	}
}

func (conf *Configuration) AddConn(conn element.Element) {

	if conf.Connectors == nil {
		conf.Connectors = make(map[string]element.Element)
		conf.Connectors [conn.Id] = conn
	} else {
		conf.Connectors[conn.Id] = conn
	}
}

func (conf *Configuration) AddAtt(a attachments.Attachment) {
	conf.Attachments = append(conf.Attachments, a)
}

func confToGoType(tConf string) string {
	foundType := false
	tGo := ""

	for t := range libraries.Repository {
		if strings.Contains(t, tConf) {
			tGo = t
			foundType = true
		}
	}

	if !foundType {
		fmt.Println("GenerateConf:: Type '" + tConf + "' NOT FOUND in Behaviour Library")
		os.Exit(0)
	}
	return tGo
}

func checkAttachments(comps map[string]ElemInfo, conns map[string]string, atts []string) {

	// Check if all components/connectors were declared
	for a := range atts {
		att := strings.Split(atts[a], ",")
		if !IsInComponents(comps, att[0]) {
			fmt.Println("Configuration:: Component '" + att[0] + "' was not Declared!!")
			os.Exit(0)
		}
		if !shared.IsInConnectors(conns, att[1]) {
			fmt.Println("Configuration:: Connector '" + att[1] + "' was not Declared!!")
			os.Exit(0)
		}
		if !IsInComponents(comps, att[2]) {
			fmt.Println("Configuration:: Component '" + att[2] + "' was not Declared!!")
			os.Exit(0)
		}
	}

	// Check if all components/connectors were used
	for c := range comps {
		if !shared.IsComponentInAttachments(atts, c) {
			fmt.Println("Configuration:: Component '" + c + "' declared, but not Used!!")
			os.Exit(0)
		}
	}

	for t := range conns {
		if !shared.IsConnectorInAttachments(atts, t) {
			fmt.Println("Configuration:: Connector '" + t + "' declared, but not Used!!")
			os.Exit(0)
		}
	}
}

type ElemInfo struct {
	ElemType string
	Param    int // TODO
}

func IsInComponents(comps map[string]ElemInfo, c string) bool {
	foundComponent := false

	for i := range comps {
		if c == i {
			foundComponent = true
			break
		}
	}
	return foundComponent
}

func MapADLIntoGo(adlFileName string) Configuration {
	conf := Configuration{ADLFileName: adlFileName}

	fullPathAdlFileName := parameters.DIR_CONF + "/" + conf.ADLFileName

	// read file
	fileContent := []string{}
	file, err := os.Open(fullPathAdlFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	// Generate Configuration
	confName := ""
	for l := range fileContent {
		if strings.Contains(strings.ToUpper(fileContent[l]), "CONFIGURATION") {
			temp := strings.Split(fileContent[l], " ")
			confName = strings.TrimSpace(temp[1])
		}
	}
	if confName == "" {
		fmt.Println("Configuration:: Configuration name not defined in '" + fullPathAdlFileName + "' ")
		os.Exit(0)
	}
	conf.Id = confName

	// Define adaptatibility
	foundAdaptability := false
	requiredAdaptations := []string{}
	for l := range fileContent {
		tempLine := fileContent[l]
		if strings.Contains(strings.ToUpper(tempLine), "ADAPTABILITY") {
			foundAdaptability = true
		} else {
			if foundAdaptability && tempLine != "" && isAdaptationType(tempLine) {
				requiredAdaptations = append(requiredAdaptations, strings.ToUpper(strings.TrimSpace(tempLine)))
			} else {
				if foundAdaptability && tempLine != "" && !isAdaptationType(tempLine) {
					break
				}
			}
		}
	}

	if !foundAdaptability || len(requiredAdaptations) == 0{
		fmt.Println("Configuration:: 'Adaptability' NOT well defined!")
		os.Exit(0)
	} else {
		for i:= range requiredAdaptations{
			switch requiredAdaptations[i] {
			case shared.EVOLUTIVE:
				parameters.IS_EVOLUTIVE = true
			case shared.CORRECTIVE:
				parameters.IS_CORRECTIVE = true
			case shared.PROACTIVE:
				parameters.IS_PROACTIVE = true
			}
		}
	}

	// Identify Components
	foundComponents := false
	comps := make(map[string]ElemInfo)
	for l := range fileContent {
		tempLine := fileContent[l]
		if strings.Contains(strings.ToUpper(tempLine), "COMPONENTS") {
			foundComponents = true
		} else {
			if foundComponents && tempLine != "" && strings.Contains(tempLine, ":") {
				temp := strings.Split(fileContent[l], ":")
				compName := strings.TrimSpace(temp[0])
				compType := ""
				if strings.Contains(tempLine, "@") {
					compType = strings.TrimSpace(temp[1][:strings.Index(temp[1], "@")])
					paramPort := strings.TrimSpace(tempLine[strings.Index(tempLine, "@")+1:])
					tempPort := parameters.SetOfPorts[paramPort]                     // TODO
					comps [compName] = ElemInfo{ElemType: compType, Param: tempPort} //TODO
				} else {
					compType = strings.TrimSpace(temp[1])
					comps [compName] = ElemInfo{ElemType: compType}
				}
			} else {
				if foundComponents && tempLine != "" && !strings.Contains(tempLine, ":") {
					break
				}
			}
		}
	}

	if len(comps) == 0 {
		fmt.Println("Configuration:: 'Components' not well formed in '" + fullPathAdlFileName + "' ")
		os.Exit(0)
	}

	// Identify Connectors
	foundConnectors := false
	conns := make(map[string]string)
	for l := range fileContent {
		tempLine := fileContent[l]
		if strings.Contains(strings.ToUpper(tempLine), "CONNECTORS") {
			foundConnectors = true
		} else {
			if foundConnectors && tempLine != "" && strings.Contains(tempLine, ":") {
				temp := strings.Split(fileContent[l], ":")
				connName := strings.TrimSpace(temp[0])
				connType := strings.TrimSpace(temp[1])
				conns [connName] = connType
			} else {
				if foundConnectors && tempLine != "" && !strings.Contains(tempLine, ":") {
					break
				}
			}
		}
	}

	if len(conns) == 0 {
		fmt.Println("Configuration:: 'Connectors' not well formed in '" + fullPathAdlFileName + "' ")
		os.Exit(0)
	}

	// Identify Attachments
	foundAttachments := false
	atts := []string{}
	for l := range fileContent {
		tempLine := fileContent[l]
		if strings.Contains(strings.ToUpper(tempLine), "ATTACHMENTS") {
			foundAttachments = true
		} else {
			if foundAttachments && tempLine != "" && strings.Contains(tempLine, ",") {
				att := strings.TrimSpace(fileContent[l])
				atts = append(atts, att)
			} else {
				if foundAttachments && tempLine != "" && !strings.Contains(tempLine, ",") {
					break
				}
			}
		}
	}

	if len(atts) == 0 {
		fmt.Println("Configuration:: 'Attachments' not well formed in '" + fullPathAdlFileName + "' ")
		os.Exit(0)
	}

	// Check attachments
	checkAttachments(comps, conns, atts)

	// Add components to configuration
	compsTemp := make(map[string]element.Element)
	for c := range comps {
		if strings.Contains(comps[c].ElemType, "SRH") {
			srhElem := srh.SRH{comps[c].Param}
			compsTemp[c] = element.Element{Id: c, TypeElem: srhElem}
		} else if strings.Contains(comps[c].ElemType, "CRH") {
			srhElem := crh.CRH{comps[c].Param}
			compsTemp[c] = element.Element{Id: c, TypeElem: srhElem}
		} else {
			compsTemp[c] = element.Element{Id: c, TypeElem: libraries.Repository[confToGoType(comps[c].ElemType)].Go}
		}
		conf.AddComp(compsTemp[c])
	}

	// add connectors to configuration
	connsTemp := make(map[string]element.Element)
	for t := range conns {
		connsTemp[t] = element.Element{Id: t, TypeElem: libraries.Repository[confToGoType(conns[t])].Go}
		conf.AddConn(connsTemp[t])
	}

	// add attachments
	for a := range atts {
		attsTemp := strings.Split(atts[a], ",")
		c1 := compsTemp[strings.TrimSpace(attsTemp[0])]
		t := connsTemp[strings.TrimSpace(attsTemp[1])]
		c2 := compsTemp[strings.TrimSpace(attsTemp[2])]
		conf.AddAtt(attachments.Attachment{C1: c1, T: t, C2: c2})
	}
	return conf
}

func isAdaptationType(line string) bool {
	r := false

	line = strings.TrimSpace(strings.ToUpper(line))
	if strings.Contains(line, shared.EVOLUTIVE) || strings.Contains(line, shared.PROACTIVE) || strings.Contains(line, shared.CORRECTIVE) {
		r = true
	}
	return r
}
