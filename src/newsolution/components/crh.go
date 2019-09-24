package components

import (
	"gmidarch/development/framework/messages"
	"strconv"
	"net"
	"gmidarch/development/artefacts/graphs"
	"encoding/binary"
	"log"
	"newsolution/element"
)

type CRH struct {
	CSP   string
	Graph graphs.GraphExecutable
	InvP  chan messages.SAMessage
	TerP  chan messages.SAMessage
	Msg   messages.SAMessage
}

func NewCRH(invP *chan messages.SAMessage, terP *chan messages.SAMessage) CRH {

	// create a new instance of Server
	r := new(CRH)

	// configure the new instance
	r.CSP = "B = InvP -> I_Process -> TerP -> B"
	r.InvP = *invP
	r.TerP = *terP
	r.Msg = messages.SAMessage{}

	// configure the state machine
	r.Graph = *graphs.NewGraph(3)
	actionChannel := make(chan messages.SAMessage)

	newEdgeInfo := graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvP, ActionType: 2, ActionChannel: &r.InvP, Message: &r.Msg}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: CRH{}.I_Process, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: &r.TerP, Message: &r.Msg}
	r.Graph.AddEdge(2, 0, newEdgeInfo)

	return *r
}

func (CRH) I_Process(msg *messages.SAMessage) {

	// check message
	args := msg.Payload.([]interface{})
	host := args[0].(string)
	port := args[1].(int)
	msgToServer := args[2].([]byte)

	// connect to server
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", host+":"+strconv.Itoa(int(port)))
		if err == nil {
			break
		}

	}

	defer conn.Close()

	// send message's size
	sizeMsgToServer := make([]byte, 4)
	l := uint32(len(msgToServer))
	binary.LittleEndian.PutUint32(sizeMsgToServer, l)
	conn.Write(sizeMsgToServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// send message
	_, err = conn.Write(msgToServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// receive message's size
	sizeMsgFromServer := make([]byte, 4)
	_, err = conn.Read(sizeMsgFromServer)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
	sizeFromServerInt := binary.LittleEndian.Uint32(sizeMsgFromServer)

	// receive reply
	msgFromServer := make([]byte, sizeFromServerInt)
	_, err = conn.Read(msgFromServer)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}
