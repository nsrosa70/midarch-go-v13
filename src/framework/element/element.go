package element

import (
	"framework/message"
	"strings"
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
		if auxCases[chosen][:2] == "I_" { // Update 'message' of sent actions
			internalAction := auxCases[chosen][:strings.LastIndex(auxCases[chosen],"_")] // TODO improve
			shared.Invoke(elem,internalAction,&msg)
			for c := range cases {
				if cases[c].Dir == reflect.SelectSend {
					cases[c].Send = reflect.ValueOf(msg)
				}
			}
		}
	}
}
