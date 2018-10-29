package main

import (
	"executionenvironment/executionenvironment"
	"shared/factories"
	"fmt"
)

func main(){

	// Start configuration
	executionenvironment.ExecutionEnvironment{}.Deploy("MiddlewareQueueingClient.conf")

	// Obtaing proxy to queueing service
	queueing := factories.FactoryQueueing()

	queueing.Subscribe("topic")

	fmt.Scanln()
}