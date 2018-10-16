package element

import (
	"framework/message"
	"reflect"
	"shared/shared"
	"graph/execgraph"
)

type Element struct {
	Id           string
	TypeElem     interface{}
	CSP string
	StateMachine execgraph.Graph
}

func (Element) Loop(elem interface{}, cases []reflect.SelectCase, auxCases []string){
	var msg message.Message

	for {
		chosen, value, _ := reflect.Select(cases)
		if cases[chosen].Dir == reflect.SelectRecv {
			msg = value.Interface().(message.Message)
		}
		if auxCases[chosen][:2] == shared.PREFIX_INTERNAL_ACTION { // Update 'message' of sent actions
			shared.Invoke(elem,auxCases[chosen],&msg)
			for c := range cases {
				if cases[c].Dir == reflect.SelectSend {
					cases[c].Send = reflect.ValueOf(msg)
				}
			}
		}
	}
}
