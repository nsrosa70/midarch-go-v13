package components

import (
	"gmidarch/shared/shared"
	"gmidarch/shared/parameters"
	"gmidarch/development/framework/messages"
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
