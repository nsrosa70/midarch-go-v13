package artefacts

import (
	"framework/element"
	"framework/configuration/attachments"
	"shared/shared"
	"strings"
	"errors"
	"strconv"
	"shared/parameters"
)

type MADL struct {
	SourceMADL        MADLFile
	ConfigurationName string
	Components        []element.ElementMADL
	Connectors        []element.ElementMADL
	Attachments       [] attachments.AttachmentMADL
	Adaptability      []string
}

func (m *MADL) Create(madlFile MADLFile) error {
	r1 := *new(error)

	m.SourceMADL = madlFile

	// Configuration
	configurationName, err := m.IdentifyConfigurationName()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}
	m.ConfigurationName = configurationName

	// Components
	components, err := m.IdentifyComponents()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}
	m.Components = components

	// Connectors
	connectors, err := m.IdentifyConnectors()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}
	m.Connectors = connectors

	// Attachments
	attachments, err := m.IdentifyAttachments()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}
	m.Attachments = attachments
	m.SetAttachmentTypes()

	adaptability, err := m.IdentifyAdaptability()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}
	m.Adaptability = adaptability

	err = m.Check()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}

	return r1
}

func (m MADL) IdentifyConfigurationName() (string, error) {
	r1 := ""
	r2 := *new(error)

	for l := range m.SourceMADL.Content {
		tempContent := m.SourceMADL.Content[l]
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

func (m MADL) IdentifyComponents() ([]element.ElementMADL, error) {
	foundComponents := false
	r1 := []element.ElementMADL{}
	r2 := *new(error)

	for l := range m.SourceMADL.Content {
		tempLine := m.SourceMADL.Content[l]
		if strings.Contains(strings.ToUpper(tempLine), "COMPONENTS") {
			foundComponents = true
		} else {
			if foundComponents && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ":") {
				temp := strings.Split(tempLine, ":")
				compId := strings.TrimSpace(temp[0])
				compType := ""
				compType = strings.TrimSpace(temp[1])
				r1 = append(r1, element.ElementMADL{ElemId: compId, ElemType: compType})
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

func (m MADL) IdentifyConnectors() ([]element.ElementMADL, error) {
	foundConnectors := false
	r1 := []element.ElementMADL{}
	r2 := *new(error)

	for l := range m.SourceMADL.Content {
		tempLine := m.SourceMADL.Content[l]
		if strings.Contains(strings.ToUpper(tempLine), "CONNECTORS") {
			foundConnectors = true
		} else {
			if foundConnectors && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ":") {
				temp := strings.Split(tempLine, ":")
				connId := strings.TrimSpace(temp[0])
				connType := strings.TrimSpace(temp[1])
				r1 = append(r1, element.ElementMADL{ElemId: connId, ElemType: connType})
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

func (m MADL) IdentifyAttachments() ([]attachments.AttachmentMADL, error) {
	r1 := []attachments.AttachmentMADL{}
	r2 := *new(error)

	// Identify Attachments
	foundAttachments := false
	for l := range m.SourceMADL.Content {
		tempLine := m.SourceMADL.Content[l]
		if strings.Contains(strings.ToUpper(tempLine), "ATTACHMENTS") {
			foundAttachments = true
		} else {
			if foundAttachments && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ",") {
				atts := strings.Split(strings.TrimSpace(tempLine), ",")
				c1Temp := atts[0]
				tTemp := atts[1]
				c2Temp := atts[2]

				c1 := element.ElementMADL{ElemId: c1Temp}
				t := element.ElementMADL{ElemId: tTemp}
				c2 := element.ElementMADL{ElemId: c2Temp}

				att := attachments.AttachmentMADL{c1, t, c2}
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

func (m MADL) IdentifyAdaptability() ([]string, error) {
	r1 := []string{}
	r2 := *new(error)

	foundAdaptability := false
	for l := range m.SourceMADL.Content {
		tempLine := m.SourceMADL.Content[l]
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

func (m *MADL) SetAttachmentTypes() {
	for a := range m.Attachments {
		c1 := m.Attachments[a].C1.ElemId
		c1Type, err := m.GetType(c1)
		shared.CheckError(err, "MADL")
		t := m.Attachments[a].T.ElemId
		tType, err := m.GetType(t)
		shared.CheckError(err, "MADL")
		c2 := m.Attachments[a].C2.ElemId
		c2Type, err := m.GetType(c2)
		shared.CheckError(err, "MADL")

		tempAttachment := attachments.AttachmentMADL{element.ElementMADL{ElemId: c1, ElemType: c1Type}, element.ElementMADL{ElemId: t, ElemType: tType}, element.ElementMADL{c2, c2Type}}
		m.Attachments[a] = tempAttachment
	}
}

func (m MADL) GetType(elemId string) (string, error) {
	r1 := ""
	r2 := *new(error)
	found := false

	for c := range m.Components {
		if m.Components[c].ElemId == elemId {
			found = true
			r1 = m.Components[c].ElemType
			break
		}
	}

	if !found {
		for c := range m.Connectors {
			if m.Connectors[c].ElemId == elemId {
				found = true
				r1 = m.Connectors[c].ElemType
				break
			}
		}
	}

	if !found {
		r2 = errors.New("MADL:: Type of element '" + elemId + "' not found.")
	}
	return r1, r2
}

func (m MADL) CreateEE() (MADL, error) {
	r1 := MADL{}
	r2 := *new(error)

	// configuration
	r1.ConfigurationName = m.ConfigurationName + "_EE"

	// Components
	comps := []element.ElementMADL{}

	comps = append(comps, element.ElementMADL{"ee", "ExecutionEnvironment"})
	comps = append(comps, element.ElementMADL{"evolutiveMonitor", "MAPEKMonitorEvolutive"}) //TODO
	comps = append(comps, element.ElementMADL{"monitor", "MAPEKMonitor"})
	comps = append(comps, element.ElementMADL{"analyser", "MAPEKAnalyser"})
	comps = append(comps, element.ElementMADL{"planner", "MAPEKPlanner"})
	comps = append(comps, element.ElementMADL{"executor", "MAPEKExecutor"})

	units := []string{}
	for i := 0; i < len(m.Components)+len(m.Connectors); i++ {
		units = append(units, "unit"+strconv.Itoa(i+1))
	}
	for i := 0; i < len(units); i++ {
		comps = append(comps, element.ElementMADL{units[i], "ExecutionUnit"})
	}

	// Connectors
	conns := [] element.ElementMADL{}

	conns = append(conns, element.ElementMADL{"t1", "OneToN"})
	conns = append(conns, element.ElementMADL{"t2", "OneWay"})
	conns = append(conns, element.ElementMADL{"t3", "OneWay"})
	conns = append(conns, element.ElementMADL{"t4", "OneWay"})
	conns = append(conns, element.ElementMADL{"t5", "OneWay"})
	conns = append(conns, element.ElementMADL{"t6", "OneWay"})

	// Attachments
	atts := []attachments.AttachmentMADL{}

	for i := 0; i < len(units); i++ {
		attC1 := element.ElementMADL{"ee", "ExecutionEnvironment"}
		attT := element.ElementMADL{"t1", "OneWay"}
		attC2 := element.ElementMADL{units[i], "ExecutionUnit"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})
	}

	attC1 := element.ElementMADL{"evolutiveMonitor", "MAPEKMonitorEvolutive"}
	attT := element.ElementMADL{"t2", "OneWay"}
	attC2 := element.ElementMADL{"monitor", "MAPEKMonitor"}
	atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

	attC1 = element.ElementMADL{"monitor", "MAPEKMonitor"}
	attT = element.ElementMADL{"t3", "OneWay"}
	attC2 = element.ElementMADL{"analyser", "MAPEKAnalyser"}
	atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

	attC1 = element.ElementMADL{"analyser", "MAPEKAnalyser"}
	attT = element.ElementMADL{"t4", "OneWay"}
	attC2 = element.ElementMADL{"planner", "MAPEKPlanner"}
	atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

	attC1 = element.ElementMADL{"planner", "MAPEKPlanner"}
	attT = element.ElementMADL{"t5", "OneWay"}
	attC2 = element.ElementMADL{"executor", "MAPEKExecutor"}
	atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

	attC1 = element.ElementMADL{"executor", "MAPEKExecutor"}
	attT = element.ElementMADL{"t6", "OneWay"}
	attC2 = element.ElementMADL{"ee", "ExecutionEnvironment"}
	atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

	// Adaptability
	adaptability := []string{}
	adaptability = append(adaptability, "None")

	// configure MADL EE
	r1.SourceMADL.FileName = strings.Replace(m.SourceMADL.FileName, parameters.MADL_EXTENSION, "", 99) + "_EE" + parameters.MADL_EXTENSION
	r1.SourceMADL.FilePath = m.SourceMADL.FilePath
	r1.Components = comps
	r1.Connectors = conns
	r1.Attachments = atts
	r1.Adaptability = adaptability

	return r1, r2
}

func (m MADL) Check() (error) {
	r1 := *new(error)

	// Check if all components/connectors were declared
	for a := range m.Attachments {

		if !m.IsInComponents(m.Attachments[a].C1) {
			r1 = errors.New("Component '" + m.Attachments[a].C1.ElemId + "' was not Declared!!")
			return r1
		}

		if !m.IsInConnectors(m.Attachments[a].T) {
			r1 = errors.New("Connector '" + m.Attachments[a].T.ElemId + "' was not Declared!!")
			return r1

		}
		if !m.IsInComponents(m.Attachments[a].C2) {
			r1 = errors.New("Component '" + m.Attachments[a].C2.ElemId + "' was not Declared!!")
			return r1
		}
	}

	// Check if all components/connectors were used
	for c := range m.Components {
		if !m.IsComponentInAttachments(m.Components[c]) {
			r1 = errors.New("Component '" + m.Components[c].ElemId + "' declared, but not Used!!")
			return r1
		}
	}

	for t := range m.Connectors {
		if !m.IsConnectorInAttachments(m.Connectors[t]) {
			r1 = errors.New("Connector '" + m.Connectors[t].ElemId + "' declared, but not Used!!")
			return r1
		}
	}
	return r1
}

func (m MADL) IsInConnectors(e element.ElementMADL) bool {
	foundConnector := false

	for i := range m.Connectors {
		if e.ElemId == m.Connectors[i].ElemId {
			foundConnector = true
			break
		}
	}
	return foundConnector
}

func (m MADL) IsInComponents(e element.ElementMADL) bool {
	foundComponent := false

	for i := range m.Components {
		if e.ElemId == m.Components[i].ElemId {
			foundComponent = true
			break
		}
	}
	return foundComponent
}

func (m MADL) IsComponentInAttachments(e element.ElementMADL) bool {
	foundComponent := false

	for a := range m.Attachments {
		if (m.Attachments[a].C1.ElemId == e.ElemId || m.Attachments[a].C2.ElemId == e.ElemId) {
			foundComponent = true
		}
	}

	return foundComponent
}

func (m MADL) IsConnectorInAttachments(e element.ElementMADL) bool {
	foundComponent := false

	for a := range m.Attachments {
		if (m.Attachments[a].T.ElemId == e.ElemId) {
			foundComponent = true
		}
	}
	return foundComponent
}
