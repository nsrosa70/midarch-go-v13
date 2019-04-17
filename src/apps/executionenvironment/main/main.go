package main

import (
	"fmt"
	"ee/ee"
)

func main() {

	// start configuration
	ee.Engine{}.Deploy("ExecutionEnvironment.conf")

	fmt.Scanln()
}
