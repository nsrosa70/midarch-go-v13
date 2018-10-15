package executionunit

import (
	"framework/message"
	"framework/element"
	"fmt"
	"strings"
	"reflect"
	"framework/configuration/commands"
	"shared/shared"
	"framework/libraries"
)

type ExecutionUnit struct{}

func (ExecutionUnit) Exec(elem element.Element, execChannels map[string]chan message.Message, channs map[string]chan message.Message, elemMaps map[string]string, chanUnit chan commands.LowLevelCommand) {

	// Define channels
	elemChannels := DefineChannels(execChannels, elem.Id)
	actions := map[string][]string{}
	behaviour := libraries.Repository[reflect.TypeOf(elem.TypeElem).String()].CSP
	actions[elem.Id] = FilterActions(strings.Split(behaviour, " "))
	individualChannels := map[string]chan message.Message{}
	for a := range actions[elem.Id] {
		individualChannels[actions[elem.Id][a]] = DefineChannel(elemChannels, actions[elem.Id][a])
	}

	// Assembly cases
	var msg message.Message
	cases := make([]reflect.SelectCase, len(individualChannels))
	auxCases := []string{}
	idx := 0

	for c := range individualChannels {
		if !shared.IsToElement(c) {
			cases[idx] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(individualChannels[c]), Send: reflect.ValueOf(msg)}
		} else {
			cases[idx] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(individualChannels[c]), Send: reflect.Value{}}
		}
		auxCases = append(auxCases, c)
		idx++
	}

	// Execute the loop of each component
	for {
		//shared.Invoke(elem.TypeElem, "Loop", individualChannels) // Individual loop
		shared.Invoke(element.Element{}, "Loop", elem.TypeElem, cases, auxCases) // Generic Loop
		select {
		case cmd := <-chanUnit: // a new management action is received
			switch cmd.Cmd {
			case commands.REPLACE_COMPONENT:
				elem = cmd.Args
			case commands.STOP:
				fmt.Println("Unit:: STOP [TODO]")
			}
		default:
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
		if (a[:2] != shared.PREFIX_INTERNAL_ACTION) {
			if strings.Contains(c, a) && c[:2] != shared.PREFIX_INTERNAL_ACTION {
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
