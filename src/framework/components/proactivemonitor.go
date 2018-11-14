package components

import (
	"framework/messages"
	"shared/shared"
)

type ProactiveMonitor struct{}

func (ProactiveMonitor) I_GenerateEvent(msg *messages.SAMessage, r *bool) {
	*msg = messages.SAMessage{Payload:shared.MonitoredEvolutiveData{"Monitored PROACTIVE data"}}
}