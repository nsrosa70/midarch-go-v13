package main

import (
	EE "executionenvironment/executionenvironment"
	"fmt"
	"shared/factories"
	"strconv"
	"time"
	"framework/messages"
)

func main() {

	// QUEUEING_HOST
	// start configuration
	EE.ExecutionEnvironment{}.Deploy("MiddlewareQueueingClient.conf")

	// proxy to naming service
	queueingroxy := factories.FactoryQueueing()
	idx := 0
	for {
		msg := messages.MessageMOM{Header:messages.Header{"Header"},PayLoad:"msg ["+strconv.Itoa(idx)+"]"}
		r:= queueingroxy.Publish("Topic01",msg)
		fmt.Println("Producer:: "+msg.PayLoad+" "+strconv.FormatBool(r))
		time.Sleep(100 * time.Millisecond)
		idx++
	}
}
