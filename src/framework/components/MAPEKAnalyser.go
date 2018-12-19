package components

import (
	"framework/messages"
	"shared/shared"
	"shared/parameters"
)

type MAPEKAnalyser struct{}

func (MAPEKAnalyser) I_Analyse(msg *messages.SAMessage, info interface{}, r *bool) {

	// Information received from Monitor
	listOfNewPlugins := msg.Payload.(shared.MonitoredEvolutiveData)

	// TODO Analyse the remaining kind of monitored data
	// Analyse the Evolutive monitored data
	analysisResult := shared.AnalysisResult{}
	analysisResult.Result = listOfNewPlugins
	analysisResult.Analysis = parameters.EVOLUTIVE_CHANGE
	*msg = messages.SAMessage{Payload: analysisResult}
}
