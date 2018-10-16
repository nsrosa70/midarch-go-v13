package adaptationmanager

import (
	"framework/element"
	"framework/configuration/configuration"
	"shared/shared"
	"framework/configuration/commands"
)

type AdaptationManager struct{}

func (AdaptationManager) Exec(conf configuration.Configuration, channsUnit map[string] chan commands.LowLevelCommand) {
	chanMAReactive := make(chan shared.MonitoredCorrectiveData)
	chanMAEvolutive := make(chan shared.MonitoredEvolutiveData)
	chanAP := make(chan shared.AnalysisResult)
	chanPE := make(chan commands.Plan)

	go new(Monitor).Exec(conf,chanMAReactive,chanMAEvolutive)
	go new(Analyser).Exec(conf,chanMAReactive,chanMAEvolutive,chanAP)
	go new(Planner).Exec(conf, chanAP,chanPE)
	go new(Executor).Exec(conf,chanPE,channsUnit)
}

type Request struct {
	Op   string
	args [] element.Element
}

