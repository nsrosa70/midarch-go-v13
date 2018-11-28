package main

import (
	"fmt"
	"core/engine"
)

func main() {

	// start configuration
	engine.ExecutionEnvironment{}.Deploy("MAPEK.conf")

	fmt.Scanln()
}
