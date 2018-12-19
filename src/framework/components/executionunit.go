package components

import (
	"framework/messages"
	"framework/element"
	"shared/shared"
	"fmt"
	"framework/configuration/commands"
)

type ExecutionUnit struct {
}

var unitMsg messages.SAMessage

func (unit ExecutionUnit) I_InitialiseUnit(msg *messages.SAMessage, info interface{}, r *bool) {
	elem := info.(element.Element)
	unitMsg = messages.SAMessage{}
	for e1 := range elem.ExecGraph.Edges {
		for e2 := range elem.ExecGraph.Edges[e1] {
			elem.ExecGraph.Edges[e1][e2].Action.Message = &unitMsg
		}
	}
}

func (unit ExecutionUnit) I_AdaptUnit(msg *messages.SAMessage, info interface{}, r *bool) {

	plan := msg.Payload.(commands.Plan)
	elem := info.(element.Element)
	fmt.Printf("UNIT:: %v ***************************************** \n", plan)
	fmt.Printf("UNIT:: %v ***************************************** \n", elem)
}

func (unit ExecutionUnit) I_Execute(msg *messages.SAMessage, info interface{}, r *bool) {
	elem := info.(element.Element)
	shared.Invoke(elem, "Loop", elem, elem.ExecGraph)
	*r = true
}