package components

import (
	"framework/messages"
	"reflect"
	"fmt"
	"os"
)

type MAPEKMonitor struct{}

func (MAPEKMonitor) I_Monitor(msg *messages.SAMessage, info interface{}, r *bool) {

	// TODO Is this check necessary
	switch reflect.TypeOf(msg.Payload).String() {
	case "shared.MonitoredCorrectiveData":
	case "shared.MonitoredEvolutiveData":
	case "shared.MonitoredProactiveData":
	default:
		fmt.Println("MAPEKMonitor:: Data Monitored is Invalid")
		os.Exit(0)
	}
}
