package requestor

import (
	"framework/message"
)

type Requestor struct{}

func (Requestor) I_PosInvP(msg *message.Message) {

	header := message.RequestHeader{Magic: "MIOP"}
	body := message.RequestBody{Op: msg.Payload.(message.Invocation).Op, Args: msg.Payload.(message.Invocation).Args}
	miop := message.MIOP{RequestHeader: header, RequestBody: body}
	toCRH := message.ToCRH{Host: msg.Payload.(message.Invocation).Host, Port: msg.Payload.(message.Invocation).Port, MIOP: miop}
	*msg = message.Message{toCRH}
}

func (Requestor) I_PosTerR(msg *message.Message) {
	payload := msg.Payload.(map[string]interface{})
	reply := payload["ReplyBody"]
	msgTemp := message.Message{reply}
	*msg = msgTemp
}
