package element

import "gmidarch/development/framework/messages"

type ElementMADL struct {
	ElemId    string
	ElemType  string
}

type ElementGo struct {
	ElemId string
	ElemType interface{}
	CSP string
}

// external actions common to all components and connectors
func (ElementGo) InvP(invP *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*invP
}

func (ElementGo) InvR(invR *chan messages.SAMessage, msg *messages.SAMessage) {
	*invR <- *msg
}

func (ElementGo) TerR(terR *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*terR
}

func (ElementGo) TerP(terP *chan messages.SAMessage, msg *messages.SAMessage) {
	*terP <- *msg
}
