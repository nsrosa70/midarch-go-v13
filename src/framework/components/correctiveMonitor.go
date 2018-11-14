package components

import (
	"framework/messages"
	"shared/shared"
)

type CorrectiveMonitor struct{}

func (CorrectiveMonitor) I_GenerateEvent(msg *messages.SAMessage, r *bool) {
	*msg = messages.SAMessage{Payload:shared.MonitoredEvolutiveData{"Monitored CORRECTIVE data"}}
}