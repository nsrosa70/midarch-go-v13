package components

import (
	"newsolution/development/artefacts/exec"
	"gmidarch/development/framework/messages"
	"newsolution/shared/shared"
	"newsolution/development/miop"
	"newsolution/development/element"
)

type Requestor struct {
	CSP   string
	Graph exec.ExecGraph
}

func NewRequestor() Requestor {

	// create a new instance of Server
	r := new(Requestor)

	return *r
}

func (r *Requestor) Configure(invP, terP, invR1, terR1, invR2, terR2, invR3, terR3 *chan messages.SAMessage) {

	// configure the state machine
	r.Graph = *exec.NewExecGraph(12)

	msg := new(interface{})
	newEdgeInfo := exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)

	args1 := make([]*interface{}, 2)  // message, (host+port)
	args1[0] = new(interface{}) // message
	args1[1] = new(interface{}) // host+port
	*args1[0] = msg
	*args1[1] = make([]interface{},2)

	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_SerialiseMIOP", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Args: args1}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR1, Message: msg}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR1, Message: msg}
	r.Graph.AddEdge(3, 4, newEdgeInfo)
	actionChannel = make(chan messages.SAMessage)

	args2 := make([]*interface{}, 2)  // host, port, msg
	args2[0] = new(interface{})
	args2[1] = new(interface{})
	*args2[0] = msg
	*args2[1] = make([] interface{},3)
	args2HostPort := make([]interface{},2)
	args2HostPort[0] = "localhost"
	args2HostPort[1] = 1313
	*args2[1] = args2HostPort
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_PreparetoCRH", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Args: args2}
	r.Graph.AddEdge(4, 5, newEdgeInfo)

	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR2, Message: msg}
	r.Graph.AddEdge(5, 6, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR2, Message: msg}
	r.Graph.AddEdge(6, 7, newEdgeInfo)

	actionChannel = make(chan messages.SAMessage)
	args3 := make([]*interface{}, 1)
	args3[0] = new(interface{})
	*args3[0] = msg
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_DeserialiseMIOP", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Args: args3}
	r.Graph.AddEdge(7, 8, newEdgeInfo)

	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR3, Message: msg}
	r.Graph.AddEdge(8, 9, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR3, Message: msg}
	r.Graph.AddEdge(9, 10, newEdgeInfo)

	actionChannel = make(chan messages.SAMessage)
	args4 := make([]*interface{}, 1)
	args4[0] = new(interface{})
	*args4[0] = msg
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_PrepareToClient", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Args: args4}
	r.Graph.AddEdge(10, 11, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	r.Graph.AddEdge(11, 0, newEdgeInfo)
}

func (Requestor) I_SerialiseMIOP(msg *messages.SAMessage, args [] interface{}) { // TODO
	inv := msg.Payload.(shared.Invocation)

	// assembly packet
	reqHeader := miop.RequestHeader{Context: "TODO", RequestId: 13, ResponseExpected: true, Key: 131313, Operation: inv.Req.Op}
	reqBody := miop.RequestBody{Body: inv.Req.Args}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 1, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{ReqHeader: reqHeader, ReqBody: reqBody}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

	// store host & port
	args[0] = inv.Host
	args[1] = inv.Port

	// configure message
	argsTemp := make([]interface{}, 1)
	argsTemp[0] = miopPacket
	msgToMarhsaller := shared.Request{Op: "marshall", Args: argsTemp}

	*msg = messages.SAMessage{Payload: msgToMarhsaller}
}

func (Requestor) I_DeserialiseMIOP(msg *messages.SAMessage) {

	argsTemp := make([]interface{}, 1)
	argsTemp[0] = msg.Payload
	msgToMarhsaller := shared.Request{Op: "unmarshall", Args: argsTemp}

	*msg = messages.SAMessage{Payload: msgToMarhsaller}
}

func (Requestor) I_PrepareToCRH(msg *messages.SAMessage, args [] interface{}) {

	toCRH := make([]interface{}, 3)
	toCRH[0] = args[0] // host
	toCRH[1] = args[1] // port
	toCRH[2] = msg.Payload.([]uint8)

	*msg = messages.SAMessage{Payload: toCRH}
}

func (Requestor) I_PrepareToClient(msg *messages.SAMessage) {
	miopPacket := msg.Payload.(miop.Packet)
	operationResult := miopPacket.Bd.RepBody.OperationResult

	*msg = messages.SAMessage{Payload: operationResult}
}
