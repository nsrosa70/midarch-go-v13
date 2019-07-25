package manager

import (
	"gmidarch/development/framework/messages"
	"gmidarch/development/creator"
	"errors"
	"gmidarch/development/checker"
	"strings"
	"strconv"
	"gmidarch/shared/parameters"
	"gmidarch/development/generator"
	"gmidarch/execution"
	"gmidarch/development/artefacts/madl"
	"gmidarch/development/artefacts/csp"
	"gmidarch/development/artefacts/graphs"
)

type Manager struct {
	MadlMid               madl.MADL
	MadlMidGo             madl.MADLGo
	CSPMid                csp.CSP
	DotsMid               map[string]csp.DOT
	MadlEE                madl.MADL
	MadlEEGo              madl.MADLGo
	DotsEE                map[string]csp.DOT
	CSPEE                 csp.CSP
	SMMid                 map[string]graphs.GraphExecutable
	SMEE                  map[string]graphs.GraphExecutable
	MapsMid               map[string]string
	MapsEE                map[string]string
	StructuralChannelsMid map[string]chan messages.SAMessage
	StructuralChannelsEE  map[string]chan messages.SAMessage
}

func (m Manager) Invoke(madlFileName string) (error) {
	r1 := *new(error)

	// MADLs
	creator := creator.Creator{}
	m.MadlMidGo, m.MadlEEGo, r1 = creator.Create(madlFileName)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}

	// Create Maps
	m.MapsMid = m.CreateMaps(m.MadlMidGo)
	m.MapsEE = m.CreateMaps(m.MadlEEGo)

	// Create Strcutural Channels
	m.StructuralChannelsMid = m.CreateStructuralChannels(m.MadlMidGo)
	m.StructuralChannelsEE = m.CreateStructuralChannels(m.MadlEEGo)

	// CSP (Mid + EE)
	generator := generator.Generator{}
	m.CSPMid, r1 = generator.GenerateCSP(m.MadlMidGo,m.MapsMid)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}
	m.CSPEE, r1 = generator.GenerateCSP(m.MadlEEGo,m.MapsEE)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}

	r1 = generator.GenerateCSPFile(m.CSPMid)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}
	r1 = generator.GenerateCSPFile(m.CSPEE)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}

	// Checker
	checker := checker.Checker{}
	isOk, r1 := checker.Check(m.CSPMid)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}
	if !isOk {
		r1 = errors.New("Manager:: CSP specifications have not passed on verification")
		return r1
	}
	isOk, r1 = checker.Check(m.CSPEE)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}
	if !isOk {
		r1 = errors.New("Manager:: CSP specifications have not passed on verification")
		return r1
	}

	// Invoke FDR - TODO (after integrating with David's solution)
	r1 = checker.GenerateDotFiles(m.CSPMid)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}
	r1 = checker.GenerateDotFiles(m.CSPEE)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}

	// DOTS
	m.DotsMid, r1 = csp.DOT{}.Create(m.CSPMid)
	if r1 != nil {
		r1 = errors.New("Manager" + r1.Error())
		return r1
	}
	m.DotsEE, r1 = csp.DOT{}.Create(m.CSPEE)
	if r1 != nil {
		r1 = errors.New("Manager" + r1.Error())
	}

	// State Machines
	m.SMMid = make(map[string]graphs.GraphExecutable)
	for i := range m.DotsMid {
		m.SMMid[i], r1 = execution.Create(m.DotsMid[i],m.StructuralChannelsMid)
		if r1 != nil {
			r1 := errors.New("Manager:: " + r1.Error())
			return r1
		}
	}

	m.SMEE = make(map[string]graphs.GraphExecutable)
	for i := range m.DotsEE {
		//g := graphs.GraphExecutable{}
		m.SMEE[i], r1 = execution.Create(m.DotsEE[i],m.StructuralChannelsEE)
		if r1 != nil {
			r1 := errors.New("Manager:: " + r1.Error())
			return r1
		}
	}

	// Deploy machines
	core := execution.Core{}
	core.Deploy(m.SMEE,m.MadlEEGo)
	core.Deploy(m.SMMid,m.MadlMidGo)

	return r1
}

func (Manager) CreateMaps(madlGo madl.MADLGo) (map[string]string) {
	r1 := make(map[string]string)

	partners := make(map[string]string)
	for i := 0; i < len(madlGo.Attachments); i++ {
		c1Id := madlGo.Attachments[i].C1.ElemId
		tId := madlGo.Attachments[i].T.ElemId
		c2Id := madlGo.Attachments[i].C2.ElemId

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
				r1[i+".e"+strconv.Itoa(c)] = p[j]
				c++
			}
		}
	}
	return r1
}

func (Manager) CreateStructuralChannels(madlGo madl.MADLGo) (map[string]chan messages.SAMessage) {
	r1 := make(map[string]chan messages.SAMessage)

	// Configure structural channels
	for i := 0; i < len(madlGo.Attachments); i++ {
		c1Id := madlGo.Attachments[i].C1.ElemId
		c2Id := madlGo.Attachments[i].C2.ElemId
		tId := madlGo.Attachments[i].T.ElemId

		// c1 -> t
		key01 := c1Id + "." + parameters.INVR + "." + tId
		key02 := tId + "." + parameters.INVP + "." + c1Id
		key03 := tId + "." + parameters.TERP + "." + c1Id
		key04 := c1Id + "." + parameters.TERR + "." + tId
		r1[key01] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		r1[key02] = r1[key01]
		r1[key03] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		r1[key04] = r1[key03]

		// t -> c2
		key01 = tId + "." + parameters.INVR + "." + c2Id
		key02 = c2Id + "." + parameters.INVP + "." + tId
		key03 = c2Id + "." + parameters.TERP + "." + tId
		key04 = tId + "." + parameters.TERR + "." + c2Id
		r1[key01] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		r1[key02] = r1[key01]
		r1[key03] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		r1[key04] = r1[key03]
	}
	return r1;
}
