package executionunit

import (
	"framework/messages"
	"framework/element"
	"fmt"
	"framework/configuration/commands"
	"shared/shared"
)

type ExecutionUnit struct{}

func (ExecutionUnit) Exec(elem element.Element, managementChann chan commands.LowLevelCommand) {

	// Configure Message, i.e., to set the address of 'msg' used in the unit
	msg := messages.SAMessage{}
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
				oldElem := elem  // TODO -> Generate *.dot of the new element
				elem = cmd.Args
				elem.ExecGraph = oldElem.ExecGraph
			case commands.STOP:
				fmt.Println("Unit:: STOP [TODO]")
			}
		default:
		}
	}
}