package components

import (
	"framework/messages"
	"framework/element"
	"framework/configuration/commands"
	"shared/shared"
)

type ExecutionUnit struct {
}

var unitMsg messages.SAMessage

// Unit initialization
func (unit ExecutionUnit) I_InitialiseUnit(msg *messages.SAMessage, info *interface{}, r *bool) {
	elem := (*info).(element.Element)
	unitMsg = messages.SAMessage{}
	for e1 := range elem.ExecGraph.Edges {
		for e2 := range elem.ExecGraph.Edges[e1] {
			elem.ExecGraph.Edges[e1][e2].Action.Message = &unitMsg
		}
	}
}

// Unit execution
func (unit ExecutionUnit) I_Execute(msg *messages.SAMessage, info interface{}, r *bool) {
	elem := info.(element.Element)
	shared.Invoke(elem, "Loop", elem, elem.ExecGraph)
	*r = true
}

// Unit adaptation
func (unit ExecutionUnit) I_AdaptUnit(msg *messages.SAMessage, info *interface{}, r *bool) {
	plan := msg.Payload.(commands.Plan)
	oldElem := (*info).(element.Element)

	// Only a single element is changed per time
	switch plan.Cmds[0].Cmd {
	case commands.REPLACE_COMPONENT: // high level command
		newElement := plan.Cmds[0].Args
		if newElement.Id == oldElem.Id { // This is the right unit to be updated
			newElement.ExecGraph = oldElem.ExecGraph // TODO -> Generate *.dot of the new element
			*info = newElement
		}
	}
}
