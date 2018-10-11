package adaptationmanager

import (
	"framework/message"
	"framework/configuration/configuration"
	"framework/configuration/commands"
)

type Executor struct{}

func (Executor) Exec(conf configuration.Configuration, channs map[string]chan message.Message, elemMaps map[string]string, chanPE chan commands.Plan, channsUnit map[string]chan commands.LowLevelCommand) {

	for {
		plan := <-chanPE
		for i := range plan.Cmds {
			switch plan.Cmds[i].Cmd {
			case commands.REPLACE_COMPONENT: // high level command
				newElement := plan.Cmds[i].Args
				id := newElement.Id
				cmd := commands.LowLevelCommand{commands.REPLACE_COMPONENT, newElement}
				channsUnit[id] <- cmd
			default:
			}
		}
	}
}
