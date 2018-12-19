package components

import (
	"framework/messages"
)

type MAPEKExecutor struct{}

func (MAPEKExecutor) I_Execute(msg *messages.SAMessage, info interface{}, r *bool) {
	/*plan := msg.Payload.(commands.Plan)

	for i := range plan.Cmds {
		switch plan.Cmds[i].Cmd {
		case commands.REPLACE_COMPONENT: // high level command
			newElement := plan.Cmds[i].Args
			id := newElement.Id
			cmd := commands.LowLevelCommand{commands.REPLACE_COMPONENT, newElement}
			//channsUnit[id] <- cmd
			fmt.Printf("MAPEKExecutor:: %v %v \n", id, cmd)
		default:
		}
	}
	*/
}
