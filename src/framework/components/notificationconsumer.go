package components

import (
"framework/messages"
"fmt"
"os"
)

type NotificationConsumer struct {}


func (c NotificationConsumer) I_PosInvP(msg *messages.SAMessage){
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Publish":
		//_args := inv.Args.([]interface{})
		//_topic := _args[0].(string)
		//_ip    := _args[1].(string)
		//_r := c.Subscribe(_topic,_ip)

		//_ter := shared.QueueingTermination{_r}
		//*msg = messages.SAMessage{_ter}
	default:
		fmt.Println("NotificationConsumer:: Operation " + inv.Op + " is not implemented by NotificationConsumer")
		os.Exit(0)
	}
}


func (NotificationConsumer) Subscribe(topic string, ip string) bool {
	r := true

	if _, ok := Subscribers[topic]; !ok {
		Subscribers[topic] = []string{}
	}

	Subscribers[topic] = append(Subscribers[topic], ip)

	return r
}

