package main

import (
	EE "executionenvironment/executionenvironment"
	"shared/factories"
	"strconv"
	"fmt"
)

func main() {

	// QUEUEING_HOST
	// start configuration
	EE.ExecutionEnvironment{}.Deploy("QueueClient.conf")

	// proxy to naming service
	queueingroxy := factories.FactoryQueueing()
	idx := 0
	for {
		msg1 := "Topic01:: msg ["+strconv.Itoa(idx)+"]"
		queueingroxy.Publish("Topic01",msg1)
		fmt.Println(msg1)
		idx++
		//time.Sleep(1000 * time.Millisecond)
	}
}
