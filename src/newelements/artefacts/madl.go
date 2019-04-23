package artefacts

import (
	"framework/element"
	"framework/configuration/attachments"
	"shared/shared"
	"strings"
	"errors"
	"fmt"
	"os"
	"shared/parameters"
)

type MADL struct {
	SourceMADL        MADLFile
	ConfigurationName string
	Components        []element.Element
	Connectors        []element.Element
	Attachments       [] attachments.Attachment
	Adaptability      string
}

func (m *MADL) GenerateMADL(madlFile MADLFile) {

	m.SourceMADL = madlFile

	// Configuration
	configurationName, err := m.IdentifyConfigurationName()
	shared.CheckError(err, "MADL")
	m.ConfigurationName = configurationName

	// Components
	components, err := m.IdentifyComponents()
	shared.CheckError(err, "MADL")
	m.Components = components

	// Connectors
	connectors, err := m.IdentifyConnectors()
	shared.CheckError(err, "MADL")
	m.Connectors = connectors

	// Attachments
	attachments, err := m.IdentifyAttachments()
	shared.CheckError(err, "MADL")
	m.Attachments = attachments
	m.SetAttachmentTypes()

	adaptability, err := m.IdentifyAdaptability()
	shared.CheckError(err, "MADL")
	m.Adaptability = adaptability
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

func (m MADL) IdentifyComponents() ([]element.Element, error) {
	foundComponents := false
	r1 := []element.Element{}
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
				r1 = append(r1, element.Element{ElemId: compId, ElemType: compType})
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

func (m MADL) IdentifyConnectors() ([]element.Element, error) {
	foundConnectors := false
	r1 := []element.Element{}
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
				r1 = append(r1, element.Element{ElemId: connId, ElemType: connType})
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

func (m MADL) IdentifyAttachments() ([]attachments.Attachment, error) {
	r1 := []attachments.Attachment{}
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

				c1 := element.Element{ElemId: c1Temp}
				t := element.Element{ElemId: tTemp}
				c2 := element.Element{ElemId: c2Temp}

				att := attachments.Attachment{c1, t, c2}
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

func (m MADL) IdentifyAdaptability() (string, error) {
	r1 := ""
	r2 := *new(error)

	foundAdaptability := false
	requiredAdaptations := []string{}
	for l := range m.SourceMADL.Content {
		tempLine := m.SourceMADL.Content[l]
		if strings.Contains(strings.ToUpper(tempLine), "ADAPTABILITY") {
			foundAdaptability = true
		} else {
			if foundAdaptability && !shared.SkipLine(tempLine) && isAdaptationType(tempLine) {
				requiredAdaptations = append(requiredAdaptations, strings.ToUpper(strings.TrimSpace(tempLine)))
			} else {
				if foundAdaptability && !shared.SkipLine(tempLine) && !isAdaptationType(tempLine) {
					break
				}
			}
		}
	}

	if !foundAdaptability || len(requiredAdaptations) == 0 {
		r2 = errors.New("'Adaptability' NOT well defined!")
	} else {
		for i := range requiredAdaptations {
			switch requiredAdaptations[i] {
			case parameters.EVOLUTIVE:
				parameters.IS_EVOLUTIVE = true
			case parameters.CORRECTIVE:
				parameters.IS_CORRECTIVE = true
			case parameters.PROACTIVE:
				parameters.IS_PROACTIVE = true
			}
		}
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

		tempAttachment := attachments.Attachment{element.Element{ElemId: c1, ElemType: c1Type}, element.Element{ElemId: t, ElemType: tType}, element.Element{c2, c2Type}}
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
