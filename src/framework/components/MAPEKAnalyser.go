package components

import (
	"framework/messages"
	"shared/shared"
	"shared/parameters"
)

type MAPEKAnalyser struct{}

func (MAPEKAnalyser) I_Analyse(msg *messages.SAMessage, info interface{}, r *bool) {
	listOfNewPlugins := msg.Payload.(shared.MonitoredEvolutiveData)
	analysisResult := shared.AnalysisResult{}
	analysisResult.Result = listOfNewPlugins
	analysisResult.Analysis = parameters.EVOLUTIVE_CHANGE
	*msg = messages.SAMessage{Payload:analysisResult}
}
