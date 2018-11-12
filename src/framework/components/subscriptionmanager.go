package components

import (
	"framework/messages"
	"shared/shared"
	"fmt"
	"os"
)

type SubscriptionManager struct{}

var SubscribersSM map[string][]SubscriberRecord

type SubscriberRecord struct {
	Host string
	Port int
}

func (SM SubscriptionManager) I_PosInvP(msg *messages.SAMessage, r *bool) {
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Subscribe":

		//Prepare invocation
		_args := inv.Args.([]interface{})
		_topic := _args[0].(string)
		_ip := _args[1].(string)
		_port := int(_args[2].(float64))

		// Invoke operation
		_r := SM.Subscribe(_topic, _ip, _port)

		// Prepare termination
		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}
	case "Unsubscribe":

		// Prepapre invocation
		_args := inv.Args.([]interface{})
		_topic := _args[0].(string)
		_ip := _args[1].(string)
		_port := int(_args[2].(float64))

		// Invoke operation
		_r := SM.Unsubscribe(_topic, _ip, _port)

		// Prepare termination
		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}
	case "GetSubscribers":

		// Invoke operation
		_r := SM.GetSubscribers()

		// Prepare termination
		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}

	default:
		fmt.Println("SubscriptionManager:: Operation " + inv.Op + " is not implemented by SubscriptionManager")
		os.Exit(0)
	}
}

// Return current subscribers
func (SM SubscriptionManager) GetSubscribers() map[string][]SubscriberRecord {

	return SubscribersSM
}

func (SM SubscriptionManager) Subscribe(topic string, ip string, port int) bool {
	r := true

	// Check if the list of subscribers has already been created
	if SubscribersSM == nil {
		SubscribersSM = make(map[string][]SubscriberRecord)
	}

	// Check if the topic already exists
	if _, ok := SubscribersSM[topic]; !ok {
		SubscribersSM[topic] = []SubscriberRecord{}
	}

	// Include new subscriber
	SubscribersSM[topic] = append(SubscribersSM[topic], SubscriberRecord{Host: ip, Port: port})

	return r
}

func (SM SubscriptionManager) Unsubscribe(topic string, ip string, port int) bool {
	r := true

	// Check if the list is empty
	if SubscribersSM == nil {
		r = false
	} else {
		records := []SubscriberRecord{}
		ok := false
		if records, ok = SubscribersSM[topic]; !ok {
			r = false
		} else {
			// Remove subscriber
			for i := range records {
				if records[i].Host == ip && records[i].Port == port {
					records[i] = records[len(records)-1] // Replace it with the last one.
					records = records[:len(records)-1]
					SubscribersSM[topic] = records
				}
			}
		}
	}

	return r
}
