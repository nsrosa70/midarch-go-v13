package components

import (
	"framework/messages"
	"framework/element"
	"shared/shared"
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

func (unit ExecutionUnit) I_ExecuteUnit(msg *messages.SAMessage, info interface{}, r *bool) {
	elem := info.(element.Element)
	shared.Invoke(elem, "Loop", elem, elem.ExecGraph)
}

func (unit ExecutionUnit) I_Nothing(msg *messages.SAMessage, info interface{}, r *bool) {
}
