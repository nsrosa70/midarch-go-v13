package main

import (
	"shared/shared"
	"os"
	"executionenvironment/executionenvironment"
	"shared/conf"
	"shared/parameters"
	"framework/components/queueing/queueing"
	"strconv"
	"fmt"
)

func main() {

	// read OS arguments
	shared.ProcessOSArguments(os.Args[1:])

	// start configuration
	EE := executionenvironment.ExecutionEnvironment{}
	EE.Exec(conf.GenerateConf(parameters.DIR_CONF+"/"+"MiddlewareQueueingClient.conf"), parameters.IS_ADAPTIVE)

	// proxy to naming service
	queueingClientProxy := queueing.LocateQueueing(parameters.QUEUEING_HOST, parameters.QUEUESERVER_PORT)

	for i:= 0 ; i < 10 ; i++ {
		queueingClientProxy.Publish("Topic 1",strconv.Itoa(i))
	}

	fmt.Println(queueingClientProxy.Consume("Topic 1"))
}
