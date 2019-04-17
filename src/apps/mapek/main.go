package main

import (
	"fmt"
	"ee/ee"
)

func main() {

	// start configuration
	ee.ExecutionEnvironment{}.Deploy("MAPEK.conf")

	fmt.Scanln()
}
