package Handlers

import (
	"net"
	"shared/net"
	"strings"
	"strconv"
	"fmt"
	"shared/errors"
	"encoding/json"
	"framework/messages"
)

func Handle(chn chan interface{}){
	var conn net.Conn
	var err error
	var ln net.Listener

	port := 1313
	addr := netshared.ResolveHostIp() + ":" + strings.TrimSpace(strconv.Itoa(port)) // TODO
	ln, err = net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		myError := errors.MyError{Source: "Subscriber", Message: "Unable to listen on port " + strconv.Itoa(port)}
		myError.ERROR()
	}

	// receive data
	for {
		if ln != nil {
			conn, err = ln.Accept()
			if err != nil {
				fmt.Println(err)
				myError := errors.MyError{Source: "Subscriber", Message: "Unable to accept connections at port " + strconv.Itoa(port)}
				myError.ERROR()
			}
		}
		jsonDecoder := json.NewDecoder(conn)
		msgMOM := messages.MessageMOM{}
		err = jsonDecoder.Decode(&msgMOM)

		if err != nil {
			fmt.Println(err)
			myError := errors.MyError{Source: "SRH", Message: "Unable to read data"}
			myError.ERROR()
		}
		chn <-msgMOM.PayLoad
	}
	return

}

func Handler(chn chan interface{}) {
	go Handle(chn)
}

func GetResult(chn chan interface{}) interface{}{

	return <- chn
}