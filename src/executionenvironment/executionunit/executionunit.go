package executionunit

import (
	"framework/message"
	"framework/element"
	"fmt"
	"strings"
	"framework/configuration/commands"
	"shared/shared"
	"os"
)

type ExecutionUnit struct{}

var msg message.Message

func (ExecutionUnit) Exec(elem element.Element, managementChann chan commands.LowLevelCommand) {

	// Configure Message, i.e., to set the address of 'msg' used in the unit
	for e1 := range elem.ExecGraph.Edges {
		for e2 := range elem.ExecGraph.Edges[e1] {
			elem.ExecGraph.Edges[e1][e2].Action.Message = &msg
		}
	}

	// Execute the loop of the element
	for {
		shared.Invoke(elem, "Loop", elem, elem.ExecGraph)
		select {
		case cmd := <-managementChann: // a new management action is received
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
