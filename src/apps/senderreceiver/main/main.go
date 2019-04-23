package main

import (
	"fmt"
	"ee/ee"
)

func main() {

	// start configuration
	ee.EE{}.Deploy("SenderReceiver.madl")

	fmt.Scanln()
}
