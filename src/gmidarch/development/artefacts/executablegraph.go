package artefacts

import (
	"gmidarch/development/framework/messages"
	"fmt"
	"strings"
	"gmidarch/shared/parameters"
	"gmidarch/shared/shared"
)

type GraphExecutable struct {
	NumNodes        int
	EdgesExecutable [][]EdgeExecutable
}

type EdgeExecutableInfo struct {
	ActionName     string
	ActionChannel  *chan messages.SAMessage                                  // Channel
	Message        *messages.SAMessage                                       // Message
	ExternalAction func(*chan messages.SAMessage, *messages.SAMessage)       // External action
	InternalAction func(elem interface{}, name string, args ... interface{}) // Internal action
}

type EdgeExecutable struct {
	From int
	To   int
	Info EdgeExecutableInfo
}

func NewGraph(n int) *GraphExecutable {
	return &GraphExecutable{
		NumNodes:        n,
		EdgesExecutable: make([][]EdgeExecutable, n),
	}
}

func (g *GraphExecutable) AddEdge(u, v int, a EdgeExecutableInfo) {
	g.EdgesExecutable[u] = append(g.EdgesExecutable[u], EdgeExecutable{From: u, To: v, Info: a})
}

func (g *GraphExecutable) AdjacentEdges(u int) []EdgeExecutable {
	return g.EdgesExecutable[u]
}

func (g GraphExecutable) Create(dot DOT, sc map[string]chan messages.SAMessage) (GraphExecutable, error) {
	r1 := NewGraph(dot.Dotgraph.NumNodes)
	r2 := *new(error)

	elemId := strings.Replace(dot.SourceDotFile.FileName, parameters.DOT_EXTENSION, "", 99)
	fmt.Println(elemId)

	eActions := EdgeExecutableInfo{}
	var msg messages.SAMessage
	for e1 := range dot.Dotgraph.EdgesDot {
		for e2 := range dot.Dotgraph.EdgesDot [e1] {
			edgeTemp := dot.Dotgraph.EdgesDot[e1][e2]
			actionNameFDR := edgeTemp.Action
			actionNameExec := ""
			if strings.Contains(actionNameFDR, ".") {
				actionNameExec = actionNameFDR[:strings.Index(actionNameFDR, ".")]
			}
			if shared.IsExternal(actionNameExec) { // External action
				key := elemId + "." + actionNameFDR
				channel := sc[key]
				params := EdgeExecutableInfo{}
				switch actionNameExec {
				case parameters.INVR:
					invr := channel
					params = EdgeExecutableInfo{ExternalAction: element.Element{}.InvR, Message: &msg, ActionChannel: &invr, ActionName: actionNameExec}
				case parameters.TERR:
					terr := channel
					params = EdgeExecutableInfo{ExternalAction: element.Element{}.TerR, Message: &msg, ActionChannel: &terr, ActionName: actionNameExec}
				case parameters.INVP:
					invp := channel
					params = EdgeExecutableInfo{ExternalAction: element.Element{}.InvP, Message: &msg, ActionChannel: &invp, ActionName: actionNameExec}
				case parameters.TERP:
					terp := channel
					params = EdgeExecutableInfo{ExternalAction: element.Element{}.TerP, Message: &msg, ActionChannel: &terp, ActionName: actionNameExec}
				}
				mapType := EdgeExecutableInfo{}
				mapType = params
				eActions = mapType
			}

			if shared.IsInternal(actionNameFDR) {
				msg := messages.SAMessage{}
				channel := make(chan messages.SAMessage)
				params := EdgeExecutableInfo{InternalAction: shared.Invoke, ActionName: actionNameFDR, Message: &msg, ActionChannel: &channel}
				mapType := params
				eActions = mapType
			}
			r1.AddEdge(edgeTemp.From, edgeTemp.To, eActions)
		}
	}

	return *r1, r2
}
