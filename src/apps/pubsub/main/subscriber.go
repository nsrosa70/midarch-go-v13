package main

import (
	"executionenvironment/executionenvironment"
	"shared/factories"
	"fmt"
)

func main(){

	// Start configuration
	executionenvironment.ExecutionEnvironment{}.Deploy("QueueClient.conf")

	// Obtaing proxy to queueing service
	queueing := factories.FactoryQueueing()

	fmt.Println(queueing.Subscribe("topic"))

	fmt.Scanln()
}