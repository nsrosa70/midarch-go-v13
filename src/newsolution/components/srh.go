package components

import (
	"log"
	"encoding/binary"
	"net"
	"strconv"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"newsolution/element"
)

type SRH struct {
	CSP      string
	Graph    graphs.GraphExecutable
	InvRChan chan messages.SAMessage
	TerRChan chan messages.SAMessage
	Msg      messages.SAMessage
	Host     string
	Port     int
}

var ln net.Listener
var conn net.Conn
var err error

func NewSRH(invR *chan messages.SAMessage, terR *chan messages.SAMessage, host string, port int) SRH {

	// create a new instance of Server
	r := new(SRH)

	// configure the new instance
	r.CSP = "B = I_SetMessage -> I_Process -> B"
	r.InvRChan = *invR
	r.TerRChan = *terR
	r.Msg = messages.SAMessage{}
	r.Host = host
	r.Port = port

	// configure the state machine
	r.Graph = *graphs.NewGraph(4)
	actionChannel := make(chan messages.SAMessage)
	args := make([] *interface{}, 2)
	args[0] = new(interface{})
	args[1] = new(interface{})
	*args[0] = r.Host
	*args[1] = r.Port

	newEdgeInfo := graphs.EdgeExecutableInfo{InternalActionWithArgs: SRH{}.I_Receive, ActionType: 3, ActionChannel: &actionChannel, Message: &r.Msg, Args: args}
	r.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: &r.InvRChan, Message: &r.Msg}
	r.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: &r.TerRChan, Message: &r.Msg}
	r.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs.EdgeExecutableInfo{InternalAction: SRH{}.I_Send, ActionType: 1, ActionChannel: &actionChannel, Message: &r.Msg}
	r.Graph.AddEdge(3, 0, newEdgeInfo)

	return *r
}

func (SRH) I_Receive(msg *messages.SAMessage, args [] *interface{}) {

	hostTemp := *args[0]
	portTemp := *args[1]
	host := hostTemp.(string)
	port := portTemp.(int)

	// create listener
	ln, err = net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// accept connections
	conn, err = ln.Accept()
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// receive message's size
	size := make([]byte, 4)
	_, err = conn.Read(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
	sizeInt := binary.LittleEndian.Uint32(size)

	// receive message
	msgTemp := make([]byte, sizeInt)
	_, err = conn.Read(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgTemp}
}

func (SRH) I_Send(msg *messages.SAMessage) {

	msgTemp := msg.Payload.([]interface{})[2].([]byte)

	// send message's size
	size := make([]byte, 4)
	l := uint32(len(msgTemp))
	binary.LittleEndian.PutUint32(size, l)
	conn.Write(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// send message
	_, err = conn.Write(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// close connection
	conn.Close()
	ln.Close()
}
