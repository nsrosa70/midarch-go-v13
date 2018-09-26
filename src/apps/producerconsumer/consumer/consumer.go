package main

import (
	EE "executionenvironment/executionenvironment"
	"fmt"
	"shared/factories"
)

func main() {

	// start configuration
	EE.ExecutionEnvironment{}.Deploy("MiddlewareQueueingClient.conf")

	// proxy to queueing service
	queueingClientProxy := factories.FactoryQueueing()

	for {
		fmt.Println("Consumer:: Here")
		fmt.Println(queueingClientProxy.Consume("Topic01"))
	}
}
