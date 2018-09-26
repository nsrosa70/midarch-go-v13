package main

import (
	EE "executionenvironment/executionenvironment"
	"shared/parameters"
	"framework/components/queueing/queueing"
	"fmt"
)

func main() {

	// start configuration
	EE.ExecutionEnvironment{}.Exec("MiddlewareQueueingClient.conf")

	// proxy to naming service
	queueingClientProxy := queueing.LocateQueueing(parameters.QUEUEING_HOST)

	for {
		fmt.Println("Consumer:: Here")
		fmt.Println(queueingClientProxy.Consume("Topic01"))
	}
}
