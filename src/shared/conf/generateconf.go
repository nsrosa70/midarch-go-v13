package conf

import (
	"framework/element"
	"framework/library"
	"strings"
	"framework/configuration/attachments"
	"framework/configuration/configuration"
	"os"
	"log"
	"bufio"
	"fmt"
	"framework/components/srh"
	"framework/components/crh"
	"shared/parameters"
)

func GenerateConf(fileName string) configuration.Configuration {
	conf := configuration.Configuration{}

	// read file
	fileContent := []string{}
	file, err := os.Open(fileName)
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
		if strings.Contains(fileContent[l], "Configuration") {
			temp := strings.Split(fileContent[l], " ")
			confName = strings.TrimSpace(temp[1])
		}
	}

	if confName == "" {
		fmt.Println("Something is Wrong in 'Configuration'")
		os.Exit(0)
	}
	conf.Id = confName

	// Identify Components
	foundComponents := false
	comps := make(map[string]ElemInfo)
	for l := range fileContent {
		tempLine := fileContent[l]
		if strings.Contains(tempLine, "Connectors") {
			foundComponents = false
		}

		if (foundComponents && tempLine != "") {
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
		}

		if strings.Contains(fileContent[l], "Components") {
			foundComponents = true
		}
	}

	if len(comps) == 0 {
		fmt.Println("Something is Wrong in 'Components'")
		os.Exit(0)
	}

	// Identify Connectors
	foundConnectors := false
	conns := make(map[string]string)
	for l := range fileContent {
		if strings.Contains(fileContent[l], "Attachments") {
			foundConnectors = false
		}

		if (foundConnectors && fileContent[l] != "") {
			temp := strings.Split(fileContent[l], ":")
			connName := strings.TrimSpace(temp[0])
			connType := strings.TrimSpace(temp[1])
			conns [connName] = connType
		}

		if strings.Contains(fileContent[l], "Connectors") {
			foundConnectors = true
		}
	}
	if len(conns) == 0 {
		fmt.Println("Something is Wrong in 'Connectors'")
		os.Exit(0)
	}

	// Identify Attachments
	foundAttachments := false
	atts := []string{}
	for l := range fileContent {
		if (foundAttachments && !strings.Contains(fileContent[l], "EndConf")) {
			att := strings.TrimSpace(fileContent[l])
			atts = append(atts, att)
		}

		if strings.Contains(fileContent[l], "Attachments") {
			foundAttachments = true
		}
	}
	if len(atts) == 0 {
		fmt.Println("Something is Wrong in 'Attachments'")
		os.Exit(0)
	}

	// Check attachments
	checkAttachments(comps, conns, atts)

	// Add components to configuration
	compsTemp := make(map[string]element.Element)
	for c := range comps {
		if strings.Contains(comps[c].ElemType, "SRH")  {
			srhElem := srh.SRH{comps[c].Param}
			compsTemp[c] = element.Element{Id: c, TypeElem: srhElem}
		} else if strings.Contains(comps[c].ElemType, "CRH") {
			srhElem := crh.CRH{comps[c].Param}
			compsTemp[c] = element.Element{Id: c, TypeElem: srhElem}
		} else {
			compsTemp[c] = element.Element{Id: c, TypeElem: library.TypeLibrary[confToGoType(comps[c].ElemType)]}
		}
		conf.AddComp(compsTemp[c])
	}

	// add connectors to configuration
	connsTemp := make(map[string]element.Element)
	for t := range conns {
		connsTemp[t] = element.Element{Id: t, TypeElem: library.TypeLibrary[confToGoType(conns[t])]}
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

func confToGoType(tConf string) string {
	foundType := false
	tGo := ""

	for t := range library.BehaviourLibrary {
		if strings.Contains(t, tConf) {
			tGo = t
			foundType = true
		}
	}

	if !foundType {
		fmt.Println("Type '" + tConf + "' NOT FOUND in Library")
		os.Exit(0)
	}
	return tGo
}

func checkAttachments(comps map[string]ElemInfo, conns map[string]string, atts []string) {

	// Check if all components/connectors were declared
	for a := range atts {
		att := strings.Split(atts[a], ",")
		if !isInComponents(comps, att[0]) {
			fmt.Println("Component '" + att[0] + "' was not Declared!!")
			os.Exit(0)
		}
		if !isInConnectors(conns, att[1]) {
			fmt.Println("Connector '" + att[1] + "' was not Declared!!")
			os.Exit(0)
		}
		if !isInComponents(comps, att[2]) {
			fmt.Println("Component '" + att[2] + "' was not Declared!!")
			os.Exit(0)
		}
	}

	// Check if all components/connectors were used
	for c := range comps {
		if !isComponentInAttachments(atts, c) {
			fmt.Println("Component '" + c + "' declared, but not Used!!")
			os.Exit(0)
		}
	}

	for t := range conns {
		if !isConnectorInAttachments(atts, t) {
			fmt.Println("Connector '" + t + "' declared, but not Used!!")
			os.Exit(0)
		}
	}
}

func isConnectorInAttachments(atts []string, t string) bool {
	foundConnector := false

	for a := range atts {
		att := strings.Split(atts[a], ",")
		if (strings.TrimSpace(att[1]) == t) {
			foundConnector = true
		}
	}

	return foundConnector
}

func isComponentInAttachments(atts []string, c string) bool {
	foundComponent := false

	for a := range atts {
		att := strings.Split(atts[a], ",")
		if (strings.TrimSpace(att[0]) == c || strings.TrimSpace(att[2]) == c) {
			foundComponent = true
		}
	}

	return foundComponent
}

func isInComponents(comps map[string]ElemInfo, c string) bool {
	foundComponent := false

	for i := range comps {
		if c == i {
			foundComponent = true
			break
		}
	}
	return foundComponent
}

func isInConnectors(conns map[string]string, t string) bool {
	foundConnector := false

	for i := range conns {
		if t == i {
			foundConnector = true
			break
		}
	}
	return foundConnector
}

type ElemInfo struct {
	ElemType string
	Param    int // TODO
}
