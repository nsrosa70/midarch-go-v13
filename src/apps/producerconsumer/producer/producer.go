package main

import (
	EE "executionenvironment/executionenvironment"
	"fmt"
	"shared/factories"
	"strconv"
	"time"
	"framework/components/queueing/queueing"
)

func main() {

	// QUEUEING_HOST
	// start configuration
	EE.ExecutionEnvironment{}.Deploy("MiddlewareQueueingClient.conf")

	// proxy to naming service
	queueingClientProxy := factories.FactoryQueueing()
	idx := 0
	for {
		msg := queueing.MessageMOM{Header:"Header",PayLoad:"msg ["+strconv.Itoa(idx)+"]"}
		r:= queueingClientProxy.Publish("Topic01",msg)
		fmt.Println("Producer:: "+msg.PayLoad+" "+strconv.FormatBool(r))
		time.Sleep(100 * time.Millisecond)
		idx++
	}
}
