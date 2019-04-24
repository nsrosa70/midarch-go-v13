package components

import (
	"framework/messages"
	"time"
	"shared/parameters"
	"shared/shared"
)

type MAPEKMonitorEvolutive struct{}

var firstTime = true
var listOfOldPlugins map[string]time.Time

func (MAPEKMonitorEvolutive) I_EvolutiveMonitoring(msg *messages.SAMessage, info *interface{}, r *bool) {
	confName := (*info).(string)
	newPlugins := []string{}
	listOfNewPlugins := make(map[string]time.Time)

	if firstTime {
		firstTime = false
		listOfOldPlugins = shared.LoadPlugins(confName)
	} else {
		listOfNewPlugins = shared.LoadPlugins(confName)
		newPlugins = shared.CheckForNewPlugins(listOfOldPlugins, listOfNewPlugins)
	}

	if len(newPlugins) > 0 {
		evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
		evolutiveMonitoredData = newPlugins
		*msg = messages.SAMessage{evolutiveMonitoredData}
	}

	listOfOldPlugins = listOfNewPlugins
	time.Sleep(parameters.MONITOR_TIME * time.Second)
}

func (MAPEKMonitorEvolutive) I_HasPlugin(msg *messages.SAMessage, info interface{}, r *bool) {

	if msg.Payload != nil {
		*r = true
	}
}

func (MAPEKMonitorEvolutive) I_HasNotPlugin(msg *messages.SAMessage, info interface{}, r *bool) {

	if msg.Payload == nil {
		*r = true
	}
}
