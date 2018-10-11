package executionunit

import (
	"framework/message"
	"framework/element"
	"fmt"
	"strings"
	"reflect"
	"framework/library"
	"framework/configuration/commands"
	"shared/shared"
)

type ExecutionUnit struct{}

func (ExecutionUnit) Exec(elem element.Element, execChannels map[string]chan message.Message, channs map[string]chan message.Message, elemMaps map[string]string, chanUnit chan commands.LowLevelCommand) {

	elemChannels := DefineChannels(execChannels, elem.Id)
	actions := map[string][]string{}
	behaviour := library.Repository[reflect.TypeOf(elem.TypeElem).String()].CSP
	actions[elem.Id] = FilterActions(strings.Split(behaviour, " "))
	individualChannels := map[string]chan message.Message{}
	for a := range actions[elem.Id] {
		individualChannels[actions[elem.Id][a]] = DefineChannel(elemChannels, actions[elem.Id][a])
	}

	for {
		shared.Invoke(elem.TypeElem, "Loop", individualChannels)
		//if elem.Id == "fibonacciinvoker" {
			fmt.Println("Execution Unit "+elem.Id)
		//}
		select {
		case cmd := <-chanUnit: // a new management action is received
			switch cmd.Cmd {
			case commands.REPLACE_COMPONENT:
				fmt.Println("Unit:: Replace")
				elem = cmd.Args
			case commands.STOP:
				fmt.Println("Unit:: STOP [TODO]")
			}
		default: // no new management action received
		}
	}
}

func DefineChannels(channels map[string]chan message.Message, elem string) map[string]chan message.Message {
	r := map[string]chan message.Message{}

	for c := range channels {
		if strings.Contains(c, elem) {
			r[c] = channels[c]
		}
	}
	return r
}

func DefineChannel(channels map[string]chan message.Message, a string) chan message.Message {
	var r chan message.Message
	found := false

	for c := range channels {
		if (a[:2] != "I_") {
			if strings.Contains(c, a) && c[:2] != "I_" {
				r = channels[c]
				found = true
				break
			}
		} else {
			if strings.Contains(c, a) {
				r = channels[c]
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println("Error: channel '" + a + "' not found")
	}

	return r
}

func FilterActions(actions []string) [] string {
	r := []string{}

	for a := range actions {
		action := actions[a]
		if strings.Contains(action, "I") || strings.Contains(action, "T") { // TODO
			if strings.Contains(action, ".") {
				action = action[:strings.Index(action, ".")]
			}
			r = append(r, action)
		}
	}
	return r
}
