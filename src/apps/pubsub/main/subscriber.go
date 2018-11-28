package main

import (
	"core/engine"
	"shared/factories"
	"fmt"
)

func main() {

	// Start configuration
	engine.ExecutionEnvironment{}.Deploy("QueueClient.conf")

	// Obtaing proxy to queueing service
	queueing := factories.FactoryQueueing()
	topic01 := "Topic01"
	topic02 := "Topic02"

	handler1,_ := queueing.Subscribe(topic01)
	handler2,_ := queueing.Subscribe(topic02)

	for {
		fmt.Println(handler1.GetResult())
		fmt.Println(handler2.GetResult())
	}
	fmt.Scanln()
}