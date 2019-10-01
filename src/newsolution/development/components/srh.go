package components

import (
	"log"
	"encoding/binary"
	"net"
	"strconv"
	"gmidarch/development/framework/messages"
	"newsolution/development/element"
	"newsolution/shared/shared"
	"newsolution/development/artefacts/exec"
)

type SRH struct {
	CSP      string
	Graph    exec.ExecGraph
	Host     string
	Port     int
}

var ln net.Listener
var conn net.Conn
var err error

func NewSRH(host string, port int) SRH {

	// create a new instance of Server
	r := new(SRH)

	// configure the new instance
	r.Host = host
	r.Port = port

	return *r
}

func (s *SRH) Configure(invR,terR *chan messages.SAMessage) {

	// configure the state machine
	s.Graph = *exec.NewExecGraph(4)
	actionChannel := make(chan messages.SAMessage)
	msg := new(messages.SAMessage)

	args1 := make([]*interface{}, 2)  // host, port, msg
	args1[0] = new(interface{})
	args1[1] = new(interface{})
	*args1[0] = msg
	*args1[1] = make([] interface{},3)
	args1HostPort := make([]interface{},2)
	args1HostPort[0] = "localhost"
	args1HostPort[1] = 1313
	*args1[1] = args1HostPort

	newEdgeInfo := exec.ExecEdgeInfo{InternalAction: shared.Invoke,ActionName:"I_Receive", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args1}
	s.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionType: 2, ActionChannel: invR, Message: msg}
	s.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionType: 2, ActionChannel: terR, Message: msg}
	s.Graph.AddEdge(2, 3, newEdgeInfo)

	args2 := make([]*interface{}, 1)
	args2[0] = new(interface{})
	*args2[0] = msg
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke,ActionName:"I_Send", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args2}
	s.Graph.AddEdge(3, 0, newEdgeInfo)
}

func (SRH) I_Receive(msg *messages.SAMessage, args [] interface{}) {

	hostTemp := args[0]
	portTemp := args[1]
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
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

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
