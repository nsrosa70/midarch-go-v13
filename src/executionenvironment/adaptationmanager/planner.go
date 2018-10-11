package adaptationmanager

import (
	"framework/message"
	"reflect"
	"framework/element"
	"plugin"
	"strings"
	"shared/shared"
	"framework/configuration/configuration"
	"shared/parameters"
	"framework/configuration/commands"
	"shared/errors"
)

type Planner struct{}

func (Planner) Exec(conf configuration.Configuration, channs map[string]chan message.Message, elemMaps map[string]string, chanAP chan shared.AnalysisResult, chanPE chan commands.Plan) {

	for {
		analysisResult := <-chanAP // receive analysis from Analyser
		switch analysisResult.Analysis {
		case parameters.EVOLUTIVE_CHANGE:
			plan := generateEvolutivePlan(conf, analysisResult, channs, elemMaps)
			chanPE <- plan // send plan to Executor
		default:
		}
	}
}

func generateEvolutivePlan(conf configuration.Configuration, analysisResult shared.AnalysisResult, channs map[string]chan message.Message, elemMaps map[string]string) commands.Plan {

	// build new plan from analysis result
	plan := commands.Plan{}
	cmds := []commands.HighLevelCommand{}

	newPlugins := reflect.ValueOf(analysisResult.Result)
	for i := 0; i < newPlugins.Len(); i++ {
		pluginName := newPlugins.Index(i).String()
		fy := loadPlugin(pluginName, "GetBehaviourExp")
		fz := loadPlugin(pluginName, "GetTypeElem")

		getBehaviourExp := fy.(func() string)
		getTypeElem := fz.(func() interface{})

		idNewElement := defineOldElement(conf, getTypeElem()) // TODO This is critical and needs to be improved in the future
		newElem := element.Element{Id: idNewElement, Behaviour: element.Element{}.BehaviourExec, TypeElem: getTypeElem(), BehaviourExp: getBehaviourExp()}

		cmd := commands.HighLevelCommand{commands.REPLACE_COMPONENT, newElem}
		cmds = append(cmds, cmd)
	}
	plan.Cmds = cmds

	return plan
}

func defineOldElement(conf configuration.Configuration, newElement interface{}) string {
	found := false
	oldElementId := ""

	// TODO check compatibility of old and new elements by type
	for i := range conf.Components {
		oldElementType := reflect.TypeOf(conf.Components[i].TypeElem).String()
		oldElementType = oldElementType[strings.LastIndex(oldElementType, ".")+1 : len(oldElementType)]
		newElementType := reflect.TypeOf(newElement).String()
		newElementType = newElementType[strings.LastIndex(newElementType, ".")+1 : len(newElementType)]
		if oldElementType == newElementType {
			oldElementId = conf.Components[i].Id
			found = true
		}
	}

	if !found {
		for i := range conf.Connectors {
			oldElementType := reflect.TypeOf(conf.Connectors[i].TypeElem).String()
			oldElementType = oldElementType[strings.LastIndex(oldElementType, ".")+1 : len(oldElementType)]
			newElementType := reflect.TypeOf(newElement).String()
			newElementType = newElementType[strings.LastIndex(newElementType, ".")+1 : len(newElementType)]
			if oldElementType == newElementType {
				oldElementId = conf.Connectors[i].Id
				found = true
			}
		}
	}

	if !found {
		myError := errors.MyError{Source: "Planner", Message: "New component NOT COMPATIBLE with Old ones"}
		myError.ERROR()
	}
	return oldElementId
}

func loadPlugin(pluginName string, symbolName string) (plugin.Symbol) {

	var lib *plugin.Plugin
	var err error

	lib, err = plugin.Open(parameters.DIR_PLUGINS + "/" + pluginName)

	if err != nil {
		myError := errors.MyError{Source: "Planner", Message: "Error on open plugin " + pluginName}
		myError.ERROR()
	}

	fx, err := lib.Lookup(symbolName)
	if err != nil {
		myError := errors.MyError{Source: "Planner", Message: "Function " + symbolName + " not found in plugin"}
		myError.ERROR()
	}

	return fx
}
