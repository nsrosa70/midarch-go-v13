package components

import (
	"framework/messages"
	"shared/shared"
	"plugin"
	"shared/parameters"
	"shared/errors"
	"framework/configuration/commands"
	"reflect"
	"framework/element"
	"strings"
)

type MAPEKPlanner struct{}

func (MAPEKPlanner) I_Plan(msg *messages.SAMessage, info *interface{}, r *bool) {

	analysisResult := msg.Payload.(shared.AnalysisResult)

	// build new plan from analysis result
	plan := commands.Plan{}
	cmds := []commands.HighLevelCommand{}
	conf := *info // Configuration is the "info" of this component

	newPlugins := reflect.ValueOf(analysisResult.Result)
	for i := 0; i < newPlugins.Len(); i++ {
		pluginName := newPlugins.Index(i).String()
		fy := loadPlugin(pluginName, "GetBehaviourExp")
		fz := loadPlugin(pluginName, "GetTypeElem")

		getBehaviourExp := fy.(func() string)
		getTypeElem := fz.(func() interface{})

		idNewElement := defineOldElement(conf, getTypeElem()) // TODO This is critical and needs to be improved in the future
		newElem := element.Element{Id: idNewElement, TypeElem: getTypeElem(), CSP: getBehaviourExp()}
		cmd := commands.HighLevelCommand{commands.REPLACE_COMPONENT, newElem}
		cmds = append(cmds, cmd)
	}
	plan.Cmds = cmds

	*msg = messages.SAMessage{Payload: plan}
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

func defineOldElement(comp interface{}, newElement interface{}) string {
	found := false
	oldElementId := ""
	components := comp.(map[string]element.Element)

	// TODO check compatibility of old and new elements by type
	for i := range components {
		oldElementType := reflect.TypeOf(components[i].TypeElem).String()
		oldElementType = oldElementType[strings.LastIndex(oldElementType, ".")+1 : len(oldElementType)]
		newElementType := reflect.TypeOf(newElement).String()
		newElementType = newElementType[strings.LastIndex(newElementType, ".")+1 : len(newElementType)]
		if oldElementType == newElementType {
			oldElementId = components[i].Id
			found = true
		}
	}

	if !found {
		myError := errors.MyError{Source: "Planner", Message: "New and old components have not COMPATIBLE types"}
		myError.ERROR()
	}
	return oldElementId
}
