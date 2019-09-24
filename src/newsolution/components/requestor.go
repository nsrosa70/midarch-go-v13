package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"newsolution/element"
	"newsolution/shared"
	"newsolution/miop"
)

type Requestor struct {
	CSP   string
	Graph graphs.GraphExecutable
	InvP  chan messages.SAMessage
	TerP  chan messages.SAMessage
	InvR1 chan messages.SAMessage
	TerR1 chan messages.SAMessage
	InvR2 chan messages.SAMessage
	TerR2 chan messages.SAMessage
	Msg   messages.SAMessage
	Args [] *interface{}
}

func NewRequestor(invP, terP, invR1, terR1, invR2, terR2  *chan messages.SAMessage) Requestor {

	// create a new instance of Server
	r := new(Requestor)

	// configure the new instance
	r.CSP = "B = InvP -> I_Serialise -> InvR -> TerR -> InvR -> TerR -> TerP -> B"
	r.InvP = *invP
	r.TerP = *terP
	r.InvR1 = *invR1
	r.TerR1 = *terR1
	r.InvR2 = *invR2
	r.TerR2 = *terR2
	r.Msg = messages.SAMessage{}
	r.Args = make([] *interface{},2)  // Host & Port
	r.Args[0] = new(interface{})
	r.Args[1] = new(interface{})

	// configure the state machine
	r.Graph = *graphs.NewGraph(8)

	newEdgeInfo := graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: &r.InvP, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalActionWithArgs: Requestor{}.I_SerialiseMIOP, Message: &r.Msg, ActionType: 3, ActionChannel: &actionChannel, Args:r.Args}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: &r.InvR1, Message: &r.Msg}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: &r.TerR1, Message: &r.Msg}
	r.Graph.AddEdge(3, 4, newEdgeInfo)
	actionChannel = make(chan messages.SAMessage)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalActionWithArgs: Requestor{}.I_PrepareToCRH, Message: &r.Msg, ActionType: 3, ActionChannel: &actionChannel, Args:r.Args}
	r.Graph.AddEdge(4, 5, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: &r.InvR2, Message: &r.Msg}
	r.Graph.AddEdge(5, 6, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: &r.TerR2, Message: &r.Msg}
	r.Graph.AddEdge(6, 7, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: &r.TerP, Message: &r.Msg}
	r.Graph.AddEdge(7, 0, newEdgeInfo)
	return *r
}

func (Requestor) I_SerialiseMIOP(msg *messages.SAMessage, args [] *interface{}) {
	inv := msg.Payload.(shared.Invocation)

	// assembly packet
	reqHeader := miop.RequestHeader{Context: "TODO", RequestId: 13, ResponseExpected: true, Key: 131313, Operation: inv.Req.Op}
	reqBody := miop.RequestBody{Body: inv.Req.Args}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 1, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{ReqHeader: reqHeader, ReqBody: reqBody}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

	// store host & port
	*args[0] = inv.Host
	*args[1] = inv.Port

	// configure message
	argsTemp := make([]interface{}, 1)
	argsTemp[0] = miopPacket
	msgToMarhsaller := shared.Request{Op: "marshall", Args: argsTemp}

	*msg = messages.SAMessage{Payload: msgToMarhsaller}
}

func (Requestor) I_PrepareToCRH(msg *messages.SAMessage, args [] *interface{}) {

	toCRH := make([]interface{},3)
	toCRH[0] = *args[0]  // host
	toCRH[1] = *args[1]  // port
	toCRH[2] = msg.Payload.([]uint8)

	*msg = messages.SAMessage{Payload: toCRH}
}
