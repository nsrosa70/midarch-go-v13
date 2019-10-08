package components

import (
	"newsolution/development/artefacts/exec"
	"strings"
	"gmidarch/development/framework/messages"
	"newsolution/development/element"
	"newsolution/shared/shared"
)

type Server struct {
	CSP   string
	Graph exec.ExecGraph
}

func NewServer() Server {

	// create a new instance of Server
	r := new(Server)

	return *r
}

func (s *Server) Configure(invP, terP *chan messages.SAMessage) Server {

	// configure the state machine
	msg := new(messages.SAMessage)
	s.Graph = *exec.NewExecGraph(3)

	actionChannel := make(chan messages.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	newEdgeInfo := exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, Message: msg, ActionChannel: invP, ActionType: 2}
	s.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Process", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info: info}
	s.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, Message: msg, ActionChannel: terP, ActionType: 2}
	s.Graph.AddEdge(2, 0, newEdgeInfo)

	return *s
}

func (Server) I_Process(msg *messages.SAMessage,info [] *interface{}) {
	msgTemp := strings.ToUpper(msg.Payload.(string))
	*msg = messages.SAMessage{Payload: msgTemp}
}
