package components

import (
	"framework/messages"
	"shared/shared"
)

type EvolutiveMonitor struct{}

func (EvolutiveMonitor) I_GenerateEventEvolutive(msg *messages.SAMessage, r *bool) {
	*msg = messages.SAMessage{Payload:shared.MonitoredEvolutiveData{"Monitored EVOLUTIVE data"}}
}