package wrapper

import "framework/message"

type Wrapper struct {}

func (e Wrapper) I_PreInvR(msg *message.Message){
	e.ReceiveFromApplication(msg)
}

func (Wrapper) SendToWrapper(msg *message.Message){
}

func (Wrapper) SendToApplication(){

}

func (Wrapper) ReceiveFromApplication(msg *message.Message){
}

func (Wrapper) ReceiveFromWrapper(){

}
