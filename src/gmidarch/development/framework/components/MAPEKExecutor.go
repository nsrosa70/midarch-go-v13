package components

import (
	"fmt"
	"gmidarch/development/framework/messages"
)

type MAPEKExecutor struct{}

func (MAPEKExecutor) I_Execute(msg *messages.SAMessage, info interface{}, r *bool) {

	// The Plan genetared by the 'Planner' is passed direct to the 'Execution Unit'
	// TODO One element is changed per time in the new implementation

	fmt.Println("Executor:: In Action")
	fmt.Printf("Executor:: [%v] \n",msg)
}
