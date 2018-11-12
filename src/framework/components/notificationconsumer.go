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
	"shared/parameters"
)

type NotificationConsumer struct{}

var ActiveConsumers = map[string]bool{}

var connNC net.Conn

func (NC NotificationConsumer) I_Notify(msg *messages.SAMessage, r *bool) {
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Notify":
		// Prepare invocation
		_args := inv.Args.([]interface{})

		// Notify subscribers provided by 'Notification Engine'
		NC.NotifySubscribers(_args[0].(string), _args[1].([]SubscriberRecord))
		_r := true // TODO, check if all subscribers are actually notified

		// Prepare termination to 'Notification Engine'
		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}
	default:
		fmt.Println("NotificationConsumer:: Operation " + inv.Op + " is not implemented by NotificationConsumer")
		os.Exit(0)
	}
}

func (NotificationConsumer) NotifySubscribers(msgToBeNotified string, listOfSubscribers []SubscriberRecord) {

	// Check if 'Active Consumers' (Consumers whose connection to Handler already exists) has been created
	if ActiveConsumers == nil {
		ActiveConsumers = make(map[string]bool, parameters.MAX_NUMBER_OF_ACTIVE_CONSUMERS)
	}

	// Notify Subscribers
	for i := range listOfSubscribers {
		host := listOfSubscribers[i].Host
		port := listOfSubscribers[i].Port
		addr := strings.Join([]string{host, strconv.Itoa(port)}, ":")

		// Check if the connection with the Handler already exists
		_, ok := ActiveConsumers[addr]
		if !ok {
			ActiveConsumers[addr] = true
			connNC, err = net.Dial("tcp", addr)

			portTmp = port
			if err != nil {
				fmt.Println(err)
				myError := errors.MyError{Source: "Notification Consumer", Message: "Unable to open connection to " + host + " : " + strconv.Itoa(port)}
				myError.ERROR()
			}
		}

		// Prepare message to be sent to Handler
		msgMOM := messages.MessageMOM{Header: messages.Header{""}, PayLoad: msgToBeNotified}

		// Send message
		encoder := json.NewEncoder(connNC)
		err = encoder.Encode(msgMOM)
		if err != nil {
			fmt.Println(err)
			myError := errors.MyError{Source: "Notification Consumer", Message: "Unable to send data to " + host + ":" + strconv.Itoa(port)}
			myError.ERROR()
		}
	}

	return
}
