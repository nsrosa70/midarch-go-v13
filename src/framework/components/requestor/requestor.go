package requestor

import (
	"framework/message"
)

type Requestor struct{}

func (e Requestor) Loop(channels map[string] chan message.Message) {
	var msgInvP, msgPosTerR message.Message
	for {
		select {
		case msgInvP = <-channels["InvP"]:
		case <-channels["I_PosInvP_requestor"]:
			e.I_PosInvP(&msgInvP)
		case channels["InvR"] <- msgInvP:
		case <-channels["TerR"]:
		case msgPosTerR = <-channels["I_PosTerR_requestor"]:
			e.I_PosTerR(&msgPosTerR)
		case channels["TerP"] <- msgPosTerR:
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
