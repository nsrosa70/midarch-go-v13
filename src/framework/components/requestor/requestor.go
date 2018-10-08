package requestor

import (
	"framework/message"
)

type Requestor struct{}

func (e Requestor) Loop(InvP, I_PosInvP, InvR, TerR, I_PosTerR, TerP chan message.Message) {
	var msgInvP, msgPosTerR message.Message
	for {
		select {
		case msgInvP = <-InvP:
		case <-I_PosInvP:
			e.I_PosInvP(&msgInvP)
		case InvR <- msgInvP:
		case <-TerR:
		case msgPosTerR = <-I_PosTerR:
			e.I_PosTerR(&msgPosTerR)
		case TerP <- msgPosTerR:
		}
	}
}

func (Requestor) I_PosInvP(msg *message.Message) {
	requestHeader := message.RequestHeader{Operation: msg.Payload.(message.Invocation).Op}
	requestBody := message.RequestBody{Args: msg.Payload.(message.Invocation).Args}
	miopHeader := message.MIOPHeader{Magic: "MIOP"}
	miopBody := message.MIOPBody{RequestHeader: requestHeader, RequestBody: requestBody}
	miop := message.MIOP{Header: miopHeader, Body: miopBody}
	toCRH := message.ToCRH{Host: msg.Payload.(message.Invocation).Host, Port: msg.Payload.(message.Invocation).Port, MIOP: miop}

	*msg = message.Message{toCRH}
}

func (Requestor) I_PosTerR(msg *message.Message) {
	payload := msg.Payload.(map[string]interface{})

	miopBody := payload["Body"]
	miopReplyBody := miopBody.(map[string]interface{})
	reply := miopReplyBody["ReplyBody"]
	msgTemp := message.Message{reply}
	*msg = msgTemp
}
