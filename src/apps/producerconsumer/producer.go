package main

import (
	EE "execution/execution"
	"fmt"
	"shared/factories"
	"strconv"
	"time"
)

func main() {

	// QUEUEING_HOST
	// start configuration
	EE.ExecutionEnvironment{}.Deploy("QueueClient.conf")

	// proxy to naming service
	queueingroxy := factories.FactoryQueueing()
	idx := 0
	for {
		msg := "msg ["+strconv.Itoa(idx)+"]"
		r:= queueingroxy.Publish("Topic01",msg)
		fmt.Println("Producer:: "+msg+" "+strconv.FormatBool(r))
		time.Sleep(100 * time.Millisecond)
		idx++
	}
}
