package main

import (
	EE "executionenvironment/executionenvironment"
	"fmt"
	"shared/factories"
	"time"
)

func main() {

	// start configuration
	EE.ExecutionEnvironment{}.Deploy("MiddlewareQueueingClient.conf")

	// proxy to queueing service
	queueingClientProxy := factories.FactoryQueueing()

	for {
		fmt.Print("Consumer:: ")
		fmt.Println(queueingClientProxy.Consume("Topic01").PayLoad)
		time.Sleep(120 * time.Millisecond)
	}
}
