package Handlers

import (
	"net"
	"strconv"
	"fmt"
	"shared/errors"
	"encoding/json"
	"framework/messages"
	"shared/net"
	"strings"
)

type HandlerNotify struct {
	Host   string
	Port int
}

var handlerChan = make(chan interface{})

func (HN HandlerNotify) Start() {
	var conn net.Conn
	var err error
	var ln net.Listener

	addr := netshared.ResolveHostIp() + ":" + strings.TrimSpace(strconv.Itoa(HN.Port))
	ln, err = net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		myError := errors.MyError{Source: "HandlerNotify", Message: "Unable to listen on port " + strconv.Itoa(HN.Port)}
		myError.ERROR()
	}

	if ln != nil {
		conn, err = ln.Accept()
		if err != nil {
			fmt.Println(err)
			myError := errors.MyError{Source: "HandlerNotify", Message: "Unable to accept connections at port " + strconv.Itoa(HN.Port)}
			myError.ERROR()
		}
	}

	// receive data
	for {
		jsonDecoder := json.NewDecoder(conn)
		msgMOM := messages.MessageMOM{}
		err = jsonDecoder.Decode(&msgMOM)

		if err != nil {
			fmt.Println(err)
			myError := errors.MyError{Source: "HandlerNotify", Message: "Unable to read data"}
			myError.ERROR()
		}
		handlerChan <- msgMOM.PayLoad
	}
	return

}

func (HN HandlerNotify) StartHandler() {
	go HN.Start()
}

func (HandlerNotify) GetResult() interface{} {

	return <-handlerChan
}
