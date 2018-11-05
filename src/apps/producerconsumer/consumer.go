package main

import (
	EE "executionenvironment/executionenvironment"
	"fmt"
	"shared/factories"
	"time"
)

func main() {

	// start configuration
	// QUEUEING_HOST
	EE.ExecutionEnvironment{}.Deploy("QueueClient.conf")

	// proxy to engine service
	queueingroxy := factories.FactoryQueueing()

	for {
		fmt.Print("Consumer:: ")
		fmt.Println(queueingroxy.Consume("Topic01").PayLoad)
		time.Sleep(120 * time.Millisecond)
	}
}
