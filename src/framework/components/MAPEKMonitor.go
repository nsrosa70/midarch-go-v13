package components

import (
	"framework/messages"
	"reflect"
	"fmt"
	"os"
)

type MAPEKMonitor struct{}

func (MAPEKMonitor) I_Monitor(msg *messages.SAMessage, r *bool) {
	switch reflect.TypeOf(msg.Payload).String() {
	case "shared.MonitoredCorrectiveData":
	case "shared.MonitoredEvolutiveData":
	case "shared.MonitoredProactiveData":
	default:
		fmt.Println("MAPEMonitor:: Data Monitored is Invalid")
		os.Exit(0)
	}
}
