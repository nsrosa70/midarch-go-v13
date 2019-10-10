package creator

import (
	"strconv"
	"strings"
	"newsolution/shared/parameters"
	"newsolution/gmidarch/development/artefacts/madl"
)

type Creator struct {}

func (Creator) CreateEE(mapp madl.MADL, kindOfAdaptability []string) (madl.MADL) {
	mEE := madl.MADL{}
	isAdaptive := true

	if len(kindOfAdaptability) == 1 && kindOfAdaptability[0] == "NONE" {
		isAdaptive = false
	}

	// configuration
	mEE.Configuration = mapp.Configuration + "_EE"

	// Components
	comps := []madl.Element{}

	comps = append(comps, madl.Element{ElemId:"ee", TypeName:"ExecutionEnvironment"})

	if isAdaptive {
		comps = append(comps, madl.Element{ElemId:"monitorevolutive", TypeName:"MAPEKMonitorEvolutive"}) //TODO
		comps = append(comps, madl.Element{ElemId:"mapekmonitor", TypeName:"MAPEKMonitor"})
		comps = append(comps, madl.Element{ElemId:"analyser", TypeName:"MAPEKAnalyser"})
		comps = append(comps, madl.Element{ElemId:"planner", TypeName:"MAPEKPlanner"})
		comps = append(comps, madl.Element{ElemId:"executor", TypeName:"MAPEKExecutor"})
	}

	units := []string{}
	for i := 0; i < len(mapp.Components)+len(mapp.Connectors); i++ {
		units = append(units, "unit"+strconv.Itoa(i+1))
	}
	for i := 0; i < len(units); i++ {
		comps = append(comps, madl.Element{ElemId:units[i], TypeName:"ExecutionUnit"})
	}

	// Connectors
	conns := [] madl.Element{}

	conns = append(conns, madl.Element{ElemId:"t1", TypeName:"OneToN"})

	if isAdaptive {
		conns = append(conns, madl.Element{ElemId:"t2", TypeName:"OneWay"})
		conns = append(conns, madl.Element{ElemId:"t3", TypeName:"OneWay"})
		conns = append(conns, madl.Element{ElemId:"t4", TypeName:"OneWay"})
		conns = append(conns, madl.Element{ElemId:"t5", TypeName:"OneWay"})
		conns = append(conns, madl.Element{ElemId:"t6", TypeName:"OneWay"})
	}

	// Attachments
	atts := []madl.Attachment{}

	for i := 0; i < len(units); i++ {
		attC1 := madl.Element{ElemId:"ee", TypeName:"ExecutionEnvironment"}
		attT := madl.Element{ElemId:"t1", TypeName:"OneToN"}
		attC2 := madl.Element{ElemId:units[i], TypeName:"ExecutionUnit"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})
	}

	if isAdaptive {
		attC1 := madl.Element{ElemId:"monitorevolutive", TypeName:"MAPEKMonitorEvolutive"}
		attT := madl.Element{ElemId:"t2", TypeName:"OneWay"}
		attC2 := madl.Element{ElemId:"mapekmonitor", TypeName:"MAPEKMonitor"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId:"mapekmonitor", TypeName:"MAPEKMonitor"}
		attT = madl.Element{ElemId:"t3", TypeName:"OneWay"}
		attC2 = madl.Element{ElemId:"analyser", TypeName:"MAPEKAnalyser"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId:"analyser", TypeName:"MAPEKAnalyser"}
		attT = madl.Element{ElemId:"t4", TypeName:"OneWay"}
		attC2 = madl.Element{ElemId:"planner", TypeName:"MAPEKPlanner"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId:"planner", TypeName:"MAPEKPlanner"}
		attT = madl.Element{ElemId:"t5", TypeName:"OneWay"}
		attC2 = madl.Element{ElemId:"executor", TypeName:"MAPEKExecutor"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId:"executor", TypeName:"MAPEKExecutor"}
		attT = madl.Element{ElemId:"t6", TypeName:"OneWay"}
		attC2 = madl.Element{ElemId:"ee", TypeName:"ExecutionEnvironment"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})
	}

	// Adaptability
	adaptability := []string{}
	adaptability = append(adaptability, "NONE") // TODO

	// configure MADL EE
	mEE.File = strings.Replace(mapp.File, parameters.MADL_EXTENSION, "", 99) + "_EE" + parameters.MADL_EXTENSION
	mEE.Path = mapp.Path
	mEE.Components = comps
	mEE.Connectors = conns
	mEE.Attachments = atts
	mEE.Adaptability = adaptability

	return mEE
}
