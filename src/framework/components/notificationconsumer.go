package components

import (
	"framework/messages"
	"shared/shared"
	"fmt"
	"os"
	"strings"
	"strconv"
	"net"
	"shared/errors"
	"encoding/json"
)

type NotificationConsumer struct{}

func (NC NotificationConsumer) I_PosInvP(msg *messages.SAMessage, r *bool) {
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Notify":
		_args := inv.Args.([]interface{})
		go NC.NotifySubscribers(_args[0].(string), _args[1].([]string))
		_r := true // TODO

		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}
	default:
		fmt.Println("NotificationConsumer:: Operation " + inv.Op + " is not implemented by NotificationConsumer")
		os.Exit(0)
	}
}

func (NotificationConsumer) NotifySubscribers(msgToBeNotified string, listOfSubscribers []string) {

	// Notify Subscribers
	host := "192.168.0.15"
	port := 1313
	addr := strings.Join([]string{host, strconv.Itoa(port)}, ":")
	conn, err = net.Dial("tcp", addr)

	portTmp = port
	if err != nil {
		fmt.Println(err)
		myError := errors.MyError{Source: "Notification Consumer", Message: "Unable to open connection to " + host + " : " + strconv.Itoa(port)}
		myError.ERROR()
	}

	msgMOM := messages.MessageMOM{Header: messages.Header{""}, PayLoad: msgToBeNotified}
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(msgMOM)
	if err != nil {
		fmt.Println(err)
		myError := errors.MyError{Source: "Notification Consumer", Message: "Unable to send data to " + host + ":" + strconv.Itoa(port)}
		myError.ERROR()
	}
	return
}
