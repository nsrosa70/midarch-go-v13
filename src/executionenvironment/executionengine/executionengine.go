package executionengine

import (
	"framework/message"
	"framework/configuration/configuration"
	"framework/configuration/commands"
	"framework/library"
	"reflect"
	"framework/element"
	"executionenvironment/executionunit"
	"shared/errors"
)

type ExecutionEngine struct{}

func (e ExecutionEngine) Exec(conf configuration.Configuration, channs map[string]chan message.Message, elemMaps map[string]string, channsUnit map[string] chan commands.LowLevelCommand){

	// configure behaviour & behaviour expression
	for i := range conf.Components {
		b := library.BehaviourLibrary[reflect.TypeOf(conf.Components[i].TypeElem).String()]
		if b == "" {
			myError := errors.MyError{Source:"Execution Engine",Message:"Component '"+conf.Components[i].Id+"' does not exist in the Library"}
			myError.ERROR()
		}
		tempElem := element.Element{conf.Components[i].Id,element.Element{}.Behaviour,conf.Components[i].TypeElem,b}
		conf.Components[i] = tempElem
	}
	for i := range conf.Connectors {
		b := library.BehaviourLibrary[reflect.TypeOf(conf.Connectors[i].TypeElem).String()]
		if b == "" {
			myError := errors.MyError{Source:"Execution Engine",Message:"Connector '"+conf.Connectors[i].Id+"'does not exist in the Library"}
			myError.ERROR()
		}
		tempElem := element.Element{conf.Connectors[i].Id,element.Element{}.Behaviour,conf.Connectors[i].TypeElem, b}
		conf.Connectors[i] = tempElem
	}

	// check behaviour using FDR
	//fdr := new(fdr.FDR)   // TODO
	//ok := fdr.CheckBehaviour(conf,elemMaps)
	//if !ok{
	//	myError := errors.MyError{Source:"Execution Engine",Message:"Configuration has a problem detected by FDR4"}
	//	myError.ERROR()
	//}

	// start components
	for i := range conf.Components {
		id := conf.Components[i].Id
		unit := new(executionunit.ExecutionUnit)
		go unit.Exec(conf.Components[i], channs, elemMaps, channsUnit[id])
	}

	// start connectors
	for i := range conf.Connectors {
		id := conf.Connectors[i].Id
		unit := new(executionunit.ExecutionUnit)
		go unit.Exec(conf.Connectors[i], channs, elemMaps, channsUnit[id])
	}
}

