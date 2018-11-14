package components

import (
	"framework/messages"
	"shared/shared"
)

type MAPEKAnalyser struct{}

func (MAPEKAnalyser) I_PosInvP(msg *messages.SAMessage, r *bool) {
	*msg = messages.SAMessage{Payload:shared.AnalysisResult{Result:"Analisys Result"}}
}
