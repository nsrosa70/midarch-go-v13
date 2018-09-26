package main

import (
	"shared/shared"
	"os"
	"executionenvironment/executionenvironment"
	"shared/conf"
	"shared/parameters"
	"framework/components/queueing/queueing"
	"fmt"
)

func main() {

	// read OS arguments
	shared.ProcessOSArguments(os.Args[1:])

	// start configuration
	EE := executionenvironment.ExecutionEnvironment{}
	EE.Exec(conf.GenerateConf("MiddlewareQueueingClient.conf"), parameters.IS_ADAPTIVE)

	// proxy to naming service
	queueingClientProxy := queueing.LocateQueueing(parameters.QUEUEING_HOST)

	for {
		r:= queueingClientProxy.Publish("Topic01","msg")
		fmt.Print("Producer:: ")
		fmt.Println(r)
	}
}
