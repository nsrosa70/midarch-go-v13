package main

import (
	"fmt"
	"executionenvironment/executionenvironment"
)

func main() {

	// start configuration
	executionenvironment.ExecutionEnvironment{}.Deploy("ExecutionEnvironment.conf")

	fmt.Scanln()
}
