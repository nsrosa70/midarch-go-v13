package components

import (
	"framework/messages"
	"time"
	"shared/parameters"
	"shared/shared"
)

type MAPEKEvolutiveMonitor struct{}

var listOfOldPlugins = shared.LoadPlugins()

func (MAPEKEvolutiveMonitor) I_EvolutiveMonitoring(msg *messages.SAMessage, info interface{}, r *bool) {
	listOfNewPlugins := shared.LoadPlugins()
	newPlugins := shared.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)

	if len(newPlugins) > 0 {
		evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
		evolutiveMonitoredData = newPlugins
		*msg = messages.SAMessage{evolutiveMonitoredData}
	}

	listOfOldPlugins = listOfNewPlugins
	time.Sleep(parameters.MONITOR_TIME * time.Second)
}

func (MAPEKEvolutiveMonitor) I_HasPlugin(msg *messages.SAMessage, info interface{}, r *bool) {

	if msg.Payload != nil {
		*r = true
	}
}

func (MAPEKEvolutiveMonitor) I_HasNotPlugin(msg *messages.SAMessage, info interface{}, r *bool) {

	if msg.Payload == nil {
		*r = true
	}
}
