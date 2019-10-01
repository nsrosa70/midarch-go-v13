package exec

import (
	"newsolution/development/artefacts/dot"
	"gmidarch/development/framework/messages"
	"strings"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
	"newsolution/development/element"
)

type Exec struct{}

func (Exec) Create(id string, dot dot.DOTGraph, maps map[string]string, channels map[string]chan messages.SAMessage) (ExecGraph) {
	r1 := NewExecGraph(dot.NumNodes)

	msg := new(messages.SAMessage)
	for e1 := range dot.EdgesDot {
		for e2 := range dot.EdgesDot [e1] {
			eActions := ExecEdgeInfo{}
			edgeTemp := dot.EdgesDot[e1][e2]
			actionNameFDR := edgeTemp.Action
			actionNameExec := ""
			if strings.Contains(actionNameFDR, ".") {
				actionNameExec = actionNameFDR[:strings.Index(actionNameFDR, ".")]
			}
			if shared.IsExternal(actionNameExec) { // External action
				actionNameTemp := strings.Split(actionNameFDR,".")
				key1 := id+"."+actionNameTemp[1]
				key2 := id+"."+ actionNameTemp[0]+"."+maps[key1]
				channel, _ := channels[key2]
				params := ExecEdgeInfo{}
				switch actionNameExec {
				case parameters.INVR:
					invr := channel
					params = ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionName:"InvR", ActionType: 2, Message: msg, ActionChannel: &invr}
				case parameters.TERR:
					terr := channel
					params = ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionName:"TerR", ActionType: 2, Message: msg, ActionChannel: &terr}
				case parameters.INVP:
					invp := channel
					params = ExecEdgeInfo{ExternalAction: element.Element{}.InvP,ActionName:"InvP", ActionType: 2, Message: msg, ActionChannel: &invp}
				case parameters.TERP:
					terp := channel
					params = ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionName:"TerP", ActionType: 2, Message: msg, ActionChannel: &terp}
				}
				mapType := ExecEdgeInfo{}
				mapType = params
				eActions = mapType
			}

			if shared.IsInternal(actionNameFDR) {  // Internal action
				channel := make(chan messages.SAMessage)
				args := make([]*interface{}, 1)
				args[0] = new(interface{})
				*args[0] = msg
				params := ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: actionNameFDR, ActionType: 1, ActionChannel: &channel, Message: msg, Args: args}
				mapType := params
				eActions = mapType
			}
			r1.AddEdge(edgeTemp.From, edgeTemp.To, eActions)
		}
	}

	return *r1
}

