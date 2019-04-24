package artefacts

import (
	"framework/element"
	"framework/configuration/attachments"
	"framework/architecturallibrary"
	"errors"
)

type MADLGO struct {
	Components   []element.ElementGo
	Connectors   []element.ElementGo
	Attachments  [] attachments.AttachmentGo
	Adaptability []string
}

func (m *MADLGO) Create(madl MADL) error {
	r1 := *new(error)
	lib := architecturallibrary.ArchitecturalLibrary{}

	err := lib.Load()
	if err != nil {
		r1 = errors.New("MADLGO:: " + err.Error())
		return r1
	}

	// Components
	comps := []element.ElementGo{}
	for c := range madl.Components {
		compMadl := madl.Components[c]
		err, compRecord := lib.GetRecord(compMadl.ElemType)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		compGoTemp := element.ElementGo{ElemId: compMadl.ElemId, ElemType: compRecord.Go, CSP:compRecord.CSP}
		comps = append(comps, compGoTemp)
	}
	m.Components = comps

	// Connectors
	conns := []element.ElementGo{}
	for c := range madl.Connectors {
		connMadl := madl.Connectors[c]
		err, connRecord := lib.GetRecord(connMadl.ElemType)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		connGoTemp := element.ElementGo{ElemId: connMadl.ElemId, ElemType: connRecord.Go, CSP:connRecord.CSP}
		conns = append(comps, connGoTemp)
	}
	m.Connectors = conns

	// Attachments
	atts := []attachments.AttachmentGo{}
	for a := range madl.Attachments {
		attMadl := madl.Attachments[a]

		c1Type := attMadl.C1.ElemType
		err, c1Record := lib.GetRecord(c1Type)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		c1Go := element.ElementGo{ElemId: attMadl.C1.ElemId, ElemType: c1Record.Go,CSP:c1Record.CSP}

		tType := attMadl.T.ElemType
		err, tRecord := lib.GetRecord(tType)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		tGo := element.ElementGo{ElemId: attMadl.T.ElemId, ElemType: tRecord.Go,CSP:tRecord.CSP}

		c2Type := attMadl.C2.ElemType
		err, c2Record := lib.GetRecord(c2Type)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		c2Go := element.ElementGo{ElemId: attMadl.C2.ElemId, ElemType: c2Record.Go,CSP:c2Record.CSP}

		atts = append(atts, attachments.AttachmentGo{c1Go, tGo, c2Go})
	}
	m.Attachments = atts

	return r1
}
