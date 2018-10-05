package srh

import (
	"net"
	"strings"
	"strconv"
	"encoding/json"
	"framework/message"
	"fmt"
	"shared/errors"
	"shared/net"
)

type SRH struct {
	Port int
}

var conn net.Conn
var err error
var ln net.Listener
var serverUp = false

func (s SRH) Loop(I_PreInvR, InvR, TerR, I_PosInvR chan message.Message) {
	var msgPreInvR,msgInvR,msgTerR,msgPosInvR message.Message
	for {
		select {
		case msgPreInvR = <-I_PreInvR:
			s.I_PreInvR(&msgPreInvR)
		case InvR <- msgInvR:
		case msgTerR = <- TerR:
		case msgPosInvR = <- I_PosInvR:
			s.I_PosTerR(&msgPosInvR)
		}
	}
}

func (s SRH) I_PreInvR(msg *message.Message) {

	if !serverUp {
		addr := netshared.ResolveHostIp() + ":" + strings.TrimSpace(strconv.Itoa(s.Port))
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			fmt.Println(err)
			myError := errors.MyError{Source: "SRH", Message: "Unable to listen on port " + strconv.Itoa(s.Port)}
			myError.ERROR()
		}
		serverUp = true
	}

	if ln != nil {
		conn, err = ln.Accept()
		if err != nil {
			fmt.Println(err)
			myError := errors.MyError{Source: "SRH", Message: "Unable to accept connections at port " + strconv.Itoa(s.Port)}
			myError.ERROR()
		}
	}

	// receive data
	jsonDecoder := json.NewDecoder(conn)
	miop := message.MIOP{}
	err = jsonDecoder.Decode(&miop)

	if err != nil {
		fmt.Println(err)
		myError := errors.MyError{Source: "SRH", Message: "Unable to read data"}
		myError.ERROR()
	}
	msg.Payload = miop
}

func (SRH) I_PosTerR(msg *message.Message) {

	// send data
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(msg)
	if err != nil {
		fmt.Println(err)
		myError := errors.MyError{Source: "SRH", Message: "Unable to send data"}
		myError.ERROR()
	}
	conn.Close()
}
