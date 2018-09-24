package main

import (
	"executionenvironment/executionengine"
	"fmt"
	"apps/conf"
)

func main() {

	// start configuration
	executionengine.StartConfiguration(conf.SenderReceiveConf())

	fmt.Println("Sender Received started!!")
	fmt.Scanln()
}
