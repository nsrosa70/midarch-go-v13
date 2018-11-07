package components

import (
	"fmt"
	"os"
	"framework/messages"
)

type NotificationEngine struct{}

func (NotificationEngine) I_PosInvP(msg *messages.SAMessage, r *bool) {

	switch msg.Payload.(messages.Invocation).Op {
	case "Subscribe":
	case "Unsubscribe":
	case "Publish":
	case "Consume":
	default:
		fmt.Println("NotificationEngine:: Operation " + msg.Payload.(messages.Invocation).Op + " is not implemented by NotificationEngine")
		os.Exit(0)
	}
}

func (NotificationEngine) I_PreInvR2(msg *messages.SAMessage, r *bool) { // Subscribe Manager

	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Subscribe":
		*r = true
	case "Unsubscribe":
		*r = true
	default:
		*r = false
	}
}

func (NotificationEngine) I_PreInvR3(msg *messages.SAMessage, r *bool) { // NOTIFICATION CONSUMER
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Publish":
		*r = true
	case "Consume":
		*r = true
	default:
		*r = false
	}
}
