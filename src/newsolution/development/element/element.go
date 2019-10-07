package element

import (
	"gmidarch/development/framework/messages"
)

type Element struct{}

func (Element) InvR(invR *chan messages.SAMessage, msg *interface{}) {
	temp := *msg

	*invR <- temp.(messages.SAMessage)
}

func (Element) TerR(terR *chan messages.SAMessage, msg *interface{}) {
	*msg = <-*terR
}

func (Element) InvP(invP *chan messages.SAMessage, msg *interface{}) {
	*msg = <-*invP
}

func (Element) TerP(terP *chan messages.SAMessage, msg *interface{}) {
	temp := *msg
	*terP <- temp.(messages.SAMessage)
}
