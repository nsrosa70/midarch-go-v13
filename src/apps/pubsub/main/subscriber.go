package main

import (
	"executionenvironment/executionenvironment"
	"shared/factories"
	"fmt"
	"shared/Handlers"
)

func main() {

	// Start configuration
	executionenvironment.ExecutionEnvironment{}.Deploy("QueueClient.conf")

	// Obtaing proxy to queueing service
	queueing := factories.FactoryQueueing()

	fmt.Println(queueing.Subscribe("Topic01"))

	chn := make(chan interface{})
	Handlers.Handler(chn)

	for {
		//x := result(chn)
		x := Handlers.GetResult(chn)
		fmt.Println(x)
	}

	fmt.Scanln()
}