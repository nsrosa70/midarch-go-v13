package madl

import (
	"newsolution/shared/shared"
	"newsolution/shared/parameters"
	"errors"
	"strings"
	"os"
	"bufio"
	"fmt"
	"reflect"
	"gmidarch/development/framework/messages"
	"strconv"
	"newsolution/development/artefacts/exec"
	"newsolution/development/repositories/architectural"
	"newsolution/development/artefacts/dot"
)

type Element struct {
	Id       string         // madl file
	TypeName string         // madl file
	Type     interface{}    // repository
	CSP      string         // repository
	Info     interface{}    // TODO
	Graph    exec.ExecGraph //
}

type Attachment struct {
	C1 Element
	T  Element
	C2 Element
}

type MADL struct {
	Path          string
	File          string
	Configuration string
	Components    []Element
	Connectors    []Element
	Attachments   []Attachment
	Adaptability  []string
	Channels      map[string]chan messages.SAMessage
	Maps          map[string]string
}

/*

func (m *MADLGo) Create(madls MADL) error {
	r1 := *new(error)
	lib := repositories.ArchitecturalLibrary{}

	// Load architectural library
	err := lib.Load()
	if err != nil {
		r1 = errors.New("MADLGO:: " + err.Error())
		return r1
	}

	// Configuration
	m.ConfigurationName = madls.ConfigurationName

	// Components
	comps := []element.ElementGo{}
	for c := range madls.Components {
		compMadl := madls.Components[c]
		err, compRecord := lib.GetRecord(compMadl.ElemType)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		compGoTemp := element.ElementGo{ElemId: compMadl.ElemId, ElemType: compRecord.Go, CSP: compRecord.CSP}
		comps = append(comps, compGoTemp)
	}
	m.Components = comps

	// Connectors
	conns := []element.ElementGo{}
	for c := range madls.Connectors {
		connMadl := madls.Connectors[c]
		err, connRecord := lib.GetRecord(connMadl.ElemType)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		connGoTemp := element.ElementGo{ElemId: connMadl.ElemId, ElemType: connRecord.Go, CSP: connRecord.CSP}
		conns = append(conns, connGoTemp)
	}
	m.Connectors = conns

	// Attachments
	atts := []attachments.AttachmentGo{}
	for a := 0; a < len(madls.Attachments); a++ {
		attMadl := madls.Attachments[a]

		c1Type := attMadl.C1.ElemType
		err, c1Record := lib.GetRecord(c1Type)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		c1Go := element.ElementGo{ElemId: attMadl.C1.ElemId, ElemType: c1Record.Go, CSP: c1Record.CSP}

		tType := attMadl.T.ElemType
		err, tRecord := lib.GetRecord(tType)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		tGo := element.ElementGo{ElemId: attMadl.T.ElemId, ElemType: tRecord.Go, CSP: tRecord.CSP}

		c2Type := attMadl.C2.ElemType
		err, c2Record := lib.GetRecord(c2Type)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		c2Go := element.ElementGo{ElemId: attMadl.C2.ElemId, ElemType: c2Record.Go, CSP: c2Record.CSP}

		atts = append(atts, attachments.AttachmentGo{c1Go, tGo, c2Go})
	}
	m.Attachments = atts

	// Adaptability
	m.Adaptability = madls.Adaptability

	return r1
}

*/

func (m *MADL) Read(file string) {

	// Read meadl file
	m.readfile(file)

	// Check configuration
	err := m.Check()
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}

	// Configure Channels and Maps
	m.configureChannelsAndMaps()

	// Configure Components
	m.configureComponents()

	// Configure Connectors
	m.configureConnectors()

}

func (m *MADL) configureComponents() {
	lib := new(architectural.ArchitecturalRepository)
	lib.Load()
	for i := range m.Components {
		record, ok := lib.Library[m.Components[i].TypeName]
		if !ok {
			fmt.Println("MADL:: Component type '" + m.Components[i].TypeName + "'not in Library")
			os.Exit(0)
		}
		m.Components[i].Type = record.Type
		dotgraph := dot.DOT{}.Read(m.Components[i].TypeName + parameters.DOT_EXTENSION)
		execgraph := exec.Exec{}.Create(m.Components[i].Id, dotgraph, m.Maps, m.Channels)

		m.Components[i].Graph = execgraph
	}
}

func (m *MADL) configureConnectors() {
	lib := new(architectural.ArchitecturalRepository)
	lib.Load()
	for i := range m.Connectors {
		record, ok := lib.Library[m.Connectors[i].TypeName]
		if !ok {
			fmt.Println("MADL:: Çonnector type '" + m.Connectors[i].TypeName + "'not in Library")
			os.Exit(0)
		}
		m.Connectors[i].Type = record.Type
		dotgraph := dot.DOT{}.Read(m.Connectors[i].TypeName + parameters.DOT_EXTENSION)
		execgraph := exec.Exec{}.Create(m.Connectors[i].Id, dotgraph, m.Maps, m.Channels)

		m.Connectors[i].Graph = execgraph
	}
}

func (m *MADL) readfile(file string) {
	// Check file name
	err := checkFileName(file)
	if err != nil {
		fmt.Println("MADL:: " + err.Error())
		os.Exit(0)
	}

	// Configure File & Path
	m.File = file
	m.Path = parameters.DIR_MADL
	fullPathAdlFileName := m.Path + "/" + m.File

	// Read file
	fileContent := []string{}
	fileTemp, err := os.Open(fullPathAdlFileName)
	if err != nil {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	defer fileTemp.Close()

	scanner := bufio.NewScanner(fileTemp)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	// Identify Configuration
	configurationName, err := m.identifyConfigurationName(fileContent)
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	m.Configuration = configurationName

	// Identify Components
	comps, err := m.identifyComponents(fileContent)
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	m.Components = comps

	// Identify Connectors
	connectors, err := m.identifyConnectors(fileContent)
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	m.Connectors = connectors

	fmt.Printf("MADL:: %v\n",m.Connectors[0].Type)

	// Identify Attachments
	attachments, err := m.identifyAttachments(fileContent)
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	m.Attachments = attachments
	//m.SetAttachmentTypes()

	// Identify adaptability
	adaptability, err := m.identifyAdaptability(fileContent)
	if (err != nil) {
		fmt.Println("MADL: " + err.Error())
		os.Exit(0)
	}
	m.Adaptability = adaptability

}

func (MADL) identifyConfigurationName(content []string) (string, error) {
	r1 := ""
	r2 := *new(error)

	for l := range content {
		tempContent := content[l]
		if strings.Contains(strings.ToUpper(tempContent), "CONFIGURATION") {
			temp := strings.Split(tempContent, " ")
			r1 = strings.TrimSpace(temp[1])
		}
	}
	if r1 == "" {
		r2 = errors.New("Configuration name not defined.")
	}
	return r1, r2
}

func (MADL) identifyComponents(content []string) ([]Element, error) {
	foundComponents := false
	r1 := []Element{}
	r2 := *new(error)

	for l := range content {
		tempLine := content[l]
		if strings.Contains(strings.ToUpper(tempLine), "COMPONENTS") {
			foundComponents = true
		} else {
			if foundComponents && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ":") {
				temp := strings.Split(tempLine, ":")
				compId := strings.TrimSpace(temp[0])
				compType := ""
				compType = strings.TrimSpace(temp[1])
				r1 = append(r1, Element{Id: compId, TypeName: compType})
			} else {
				if foundComponents && !shared.SkipLine(tempLine) && !strings.Contains(tempLine, ":") {
					break
				}
			}
		}
	}

	if len(r1) == 0 {
		r2 = errors.New("MADL:: 'Components' not well formed.")
	}

	return r1, r2
}

func (MADL) identifyConnectors(content []string) ([]Element, error) {
	foundConnectors := false
	r1 := []Element{}
	r2 := *new(error)

	for l := range content {
		tempLine := content[l]
		if strings.Contains(strings.ToUpper(tempLine), "CONNECTORS") {
			foundConnectors = true
		} else {
			if foundConnectors && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ":") {
				temp := strings.Split(tempLine, ":")
				connId := strings.TrimSpace(temp[0])
				connType := strings.TrimSpace(temp[1])
				connTypeName := connType
				r1 = append(r1, Element{Id: connId, Type: connType, TypeName:connTypeName})
			} else {
				if foundConnectors && tempLine != "" && !strings.Contains(tempLine, ":") {
					break
				}
			}
		}
	}

	if len(r1) == 0 {
		r2 = errors.New("MADL:: 'Connectors' not well formed.")
	}

	return r1, r2
}

func (MADL) identifyAttachments(content []string) ([]Attachment, error) {
	r1 := []Attachment{}
	r2 := *new(error)

	// Identify Attachments
	foundAttachments := false
	for l := range content {
		tempLine := content[l]
		if strings.Contains(strings.ToUpper(tempLine), "ATTACHMENTS") {
			foundAttachments = true
		} else {
			if foundAttachments && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ",") {
				atts := strings.Split(strings.TrimSpace(tempLine), ",")
				c1Temp := atts[0]
				tTemp := atts[1]
				c2Temp := atts[2]

				c1 := Element{Id: c1Temp}
				t := Element{Id: tTemp}
				c2 := Element{Id: c2Temp}

				att := Attachment{c1, t, c2}
				r1 = append(r1, att)
			} else {
				if foundAttachments && tempLine != "" && !strings.Contains(tempLine, ",") {
					break
				}
			}
		}
	}

	if len(r1) == 0 {
		r2 = errors.New("MADL:: 'Attachments' not well formed.")
	}

	return r1, r2
}

func (MADL) identifyAdaptability(content []string) ([]string, error) {
	r1 := []string{}
	r2 := *new(error)

	foundAdaptability := false
	for l := range content {
		tempLine := content[l]
		if strings.Contains(strings.ToUpper(tempLine), "ADAPTABILITY") {
			foundAdaptability = true
		} else {
			if foundAdaptability && !shared.SkipLine(tempLine) && shared.IsAdaptationType(tempLine) {
				r1 = append(r1, strings.ToUpper(strings.TrimSpace(tempLine)))
			} else {
				if foundAdaptability && !shared.SkipLine(tempLine) && !shared.IsAdaptationType(tempLine) {
					break
				}
			}
		}
	}

	if !foundAdaptability || len(r1) == 0 {
		r2 = errors.New("'Adaptability' NOT well defined!")
	}

	return r1, r2
}

func (m MADL) PrintComponents() {

	for i := range m.Components {
		fmt.Println(reflect.TypeOf(m.Components[i].Type))
	}
}

/*
func (m *MADL) SetAttachmentTypes() {

	for a := range m.Attachments {
		c1 := m.Attachments[a].C1.Id
		c1Type, err := m.GetType(c1)
		shared.CheckError(err, "MADL")
		t := m.Attachments[a].T.Id
		tType, err := m.GetType(t)
		shared.CheckError(err, "MADL")
		c2 := m.Attachments[a].C2.Id
		c2Type, err := m.GetType(c2)
		shared.CheckError(err, "MADL")

		tempAttachment := Attachment{Element{Id: c1, Type: c1Type}, Element{Id: t, Type: tType}, Element{Id:c2, Type:c2Type}}
		m.Attachments[a] = tempAttachment
	}
}
*/

/*
func (m MADL) GetType(elemId string) (string, error) {
	r1 := ""
	r2 := *new(error)
	found := false

	for c := range m.Components {
		if m.Components[c].Id == elemId {
			found = true
			r1 = m.Components[c].Type
			break
		}
	}

	if !found {
		for c := range m.Connectors {
			if m.Connectors[c].Id == elemId {
				found = true
				r1 = m.Connectors[c].Type
				break
			}
		}
	}

	if !found {
		r2 = errors.New("MADL:: Type of element '" + elemId + "' not found.")
	}
	return r1, r2
}
*/
/*
func (m MADL) CreateEE(kindOfAdaptability []string) (MADL, error) {
	r1 := MADL{}
	r2 := *new(error)
	isAdaptive := true

	if len(kindOfAdaptability) == 1 && kindOfAdaptability[0] == "NONE" {
		isAdaptive = false
	}

	// configuration
	r1.Configuration = m.Configuration + "_EE"

	// Components
	comps := []element.ElementMADL{}

	comps = append(comps, Element{"ee", "ExecutionEnvironment"})

	if isAdaptive {
		comps = append(comps, Element{"monitorevolutive", "MAPEKMonitorEvolutive"}) //TODO
		comps = append(comps, Element{"mapekmonitor", "MAPEKMonitor"})
		comps = append(comps, Element{"analyser", "MAPEKAnalyser"})
		comps = append(comps, Element{"planner", "MAPEKPlanner"})
		comps = append(comps, Element{"executor", "MAPEKExecutor"})
	}

	units := []string{}
	for i := 0; i < len(m.Components)+len(m.Connectors); i++ {
		units = append(units, "unit"+strconv.Itoa(i+1))
	}
	for i := 0; i < len(units); i++ {
		comps = append(comps, Element{units[i], "ExecutionUnit"})
	}

	// Connectors
	conns := [] element.ElementMADL{}

	conns = append(conns, Element{"t1", "OneToN"})

	if isAdaptive {
		conns = append(conns, Element{"t2", "OneWay"})
		conns = append(conns, Element{"t3", "OneWay"})
		conns = append(conns, Element{"t4", "OneWay"})
		conns = append(conns, Element{"t5", "OneWay"})
		conns = append(conns, Element{"t6", "OneWay"})
	}

	// Attachments
	atts := []attachments.AttachmentMADL{}

	for i := 0; i < len(units); i++ {
		attC1 := Element{"ee", "ExecutionEnvironment"}
		attT := Element{"t1", "OneToN"}
		attC2 := Element{units[i], "ExecutionUnit"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})
	}

	if isAdaptive {
		attC1 := Element{"monitorevolutive", "MAPEKMonitorEvolutive"}
		attT := Element{"t2", "OneWay"}
		attC2 := Element{"mapekmonitor", "MAPEKMonitor"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

		attC1 = Element{"mapekmonitor", "MAPEKMonitor"}
		attT = Element{"t3", "OneWay"}
		attC2 = Element{"analyser", "MAPEKAnalyser"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

		attC1 = Element{"analyser", "MAPEKAnalyser"}
		attT = Element{"t4", "OneWay"}
		attC2 = Element{"planner", "MAPEKPlanner"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

		attC1 = Element{"planner", "MAPEKPlanner"}
		attT = Element{"t5", "OneWay"}
		attC2 = Element{"executor", "MAPEKExecutor"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

		attC1 = Element{"executor", "MAPEKExecutor"}
		attT = Element{"t6", "OneWay"}
		attC2 = Element{"ee", "ExecutionEnvironment"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})
	}

	// Adaptability
	adaptability := []string{}
	adaptability = append(adaptability, "NONE") // TODO

	// configure MADL EE
	r1.File = strings.Replace(m.File, parameters.MADL_EXTENSION, "", 99) + "_EE" + parameters.MADL_EXTENSION
	r1.Path = m.Path
	r1.Components = comps
	r1.Connectors = conns
	r1.Attachments = atts
	r1.Adaptability = adaptability

	return r1, r2
}
*/

/*
func (m MADL) Print() {
	lines := []string{}

	// Configuration
	lines = append(lines, "Configuration "+m.Configuration+" := "+"\n")

	// Components
	lines = append(lines, "Components"+"\n")
	for i := range m.Components {
		lines = append(lines, "   "+m.Components[i].Id+" : "+m.Components[i].Type+"\n")
	}

	// Connectors
	lines = append(lines, "Connectors"+"\n")
	for i := range m.Connectors {
		lines = append(lines, "   "+m.Connectors[i].Id+" : "+m.Connectors[i].Type+"\n")
	}

	// Attachments
	lines = append(lines, "Attachments"+"\n")
	for i := range m.Attachments {
		lines = append(lines, "   "+m.Attachments[i].C1.Id+","+m.Attachments[i].T.Id+","+m.Attachments[i].C2.Id+"\n")
	}

	// Adaptability
	lines = append(lines, "Adaptability"+"\n")
	lines = append(lines, "   "+m.Adaptability[0]+"\n") // TODO

	lines = append(lines, "EndConf")

	for i := range lines {
		fmt.Print(lines[i])
	}
}
*/

func (m MADL) Check() (error) {
	r1 := *new(error)

	// Check if all components/connectors were declared
	for a := range m.Attachments {

		if !m.isInComponents(m.Attachments[a].C1) {
			r1 = errors.New("Component '" + m.Attachments[a].C1.Id + "' was not Declared!!")
			return r1
		}

		if !m.isInConnectors(m.Attachments[a].T) {
			r1 = errors.New("Connector '" + m.Attachments[a].T.Id + "' was not Declared!!")
			return r1

		}
		if !m.isInComponents(m.Attachments[a].C2) {
			r1 = errors.New("Component '" + m.Attachments[a].C2.Id + "' was not Declared!!")
			return r1
		}
	}

	// Check if all components/connectors were used
	for c := range m.Components {
		if !m.isComponentInAttachments(m.Components[c]) {
			r1 = errors.New("Component '" + m.Components[c].Id + "' declared, but not Used!!")
			return r1
		}
	}

	for t := range m.Connectors {
		if !m.isConnectorInAttachments(m.Connectors[t]) {
			r1 = errors.New("Connector '" + m.Connectors[t].Id + "' declared, but not Used!!")
			return r1
		}
	}
	return r1
}

func (m MADL) isInConnectors(e Element) bool {
	foundConnector := false

	for i := range m.Connectors {
		if e.Id == m.Connectors[i].Id {
			foundConnector = true
			break
		}
	}
	return foundConnector
}

func (m MADL) isInComponents(e Element) bool {
	foundComponent := false

	for i := range m.Components {
		if e.Id == m.Components[i].Id {
			foundComponent = true
			break
		}
	}
	return foundComponent
}

func (m MADL) isComponentInAttachments(e Element) bool {
	foundComponent := false

	for a := range m.Attachments {
		if (m.Attachments[a].C1.Id == e.Id || m.Attachments[a].C2.Id == e.Id) {
			foundComponent = true
		}
	}

	return foundComponent
}

func (m MADL) isConnectorInAttachments(e Element) bool {
	foundComponent := false

	for a := range m.Attachments {
		if (m.Attachments[a].T.Id == e.Id) {
			foundComponent = true
		}
	}
	return foundComponent
}

func checkFileName(fileName string) error {
	r := *new(error)
	r = nil

	len := len(fileName)

	if len <= 5 {
		r = errors.New("File Name Invalid")
	} else {
		if fileName[len-5:] != parameters.MADL_EXTENSION {
			r = errors.New("Invalid extension of '" + fileName + "'")
		}
	}

	return r
}

func (madl *MADL) configureChannelsAndMaps() {
	structuralChannels := make(map[string]chan messages.SAMessage)

	// Configure structural channels
	for i := range madl.Attachments {
		c1Id := madl.Attachments[i].C1.Id
		c2Id := madl.Attachments[i].C2.Id
		tId := madl.Attachments[i].T.Id

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
	madl.Channels = structuralChannels

	// Configure maps
	elemMaps := make(map[string]string)
	partners := make(map[string]string)

	for i := range madl.Attachments {
		c1Id := madl.Attachments[i].C1.Id
		c2Id := madl.Attachments[i].C2.Id
		tId := madl.Attachments[i].T.Id
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
	madl.Maps = elemMaps
}
