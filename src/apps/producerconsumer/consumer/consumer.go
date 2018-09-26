package main

import (
	EE "executionenvironment/executionenvironment"
	"framework/components/queueing/queueing"
	"fmt"
)

func main() {

	// start configuration
	EE.ExecutionEnvironment{}.Deploy("MiddlewareQueueingClient.conf")

	// proxy to naming service
	queueingClientProxy := queueing.LocateQueueing()

	for {
		fmt.Println("Consumer:: Here")
		fmt.Println(queueingClientProxy.Consume("Topic01"))
	}
}
