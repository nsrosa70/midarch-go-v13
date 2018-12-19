package commands

import "framework/element"

const REPLACE_COMPONENT = "REPLACE_COMPONENT"
const STOP              = "STOP"
const FDR_COMMAND       = "refines"

type LowLevelCommand struct {
	Cmd      string
	Args element.Element
}

type HighLevelCommand struct {
	Cmd      string
	Args     element.Element
}

type Plan struct {
	Cmds [] HighLevelCommand
}

