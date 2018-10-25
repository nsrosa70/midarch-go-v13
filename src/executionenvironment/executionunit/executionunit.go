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
	"os"
)

type ExecutionUnit struct{}

func (ExecutionUnit) Exec(elem element.Element, strcuturalChannels map[string]chan message.Message, chanUnit chan commands.LowLevelCommand) {

	// Define channels
	actions := map[string][]string{}
	individualChannels := map[string]chan message.Message{}

	elemChannels := DefineChannels(strcuturalChannels, elem.Id)
	behaviour := libraries.Repository[reflect.TypeOf(elem.TypeElem).String()].CSP
	actions[elem.Id] = FilterActions(strings.Split(behaviour, " "))

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

	// Execute the loop of the element
	for {
		//shared.Invoke(elem.TypeElem, "Loop", individualChannels) // Individual loop
		shared.Invoke(element.Element{}, "Loop", elem.TypeElem, cases, auxCases,elem.ExitPoints) // Generic Loop
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
		fmt.Println("ExecutionUnit: channel '" + a + "' not found")
		os.Exit(0)
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
