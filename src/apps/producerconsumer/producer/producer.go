package main

import (
	EE "executionenvironment/executionenvironment"
	"fmt"
	"shared/factories"
)

func main() {

	// start configuration
	EE.ExecutionEnvironment{}.Deploy("MiddlewareQueueingClient.conf")

	// proxy to naming service
	queueingClientProxy := factories.FactoryQueueing()
	for {
		r:= queueingClientProxy.Publish("Topic01","msg")
		fmt.Print("Producer:: ")
		fmt.Println(r)
	}
}
