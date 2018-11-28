package main

import (
	"fmt"
	"core/engine"
)

func main() {

	// start configuration
	engine.Engine{}.Deploy("ExecutionEnvironment.conf")

	fmt.Scanln()
}
