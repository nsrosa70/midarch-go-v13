package main

import (
	"fmt"
	"shared/parameters"
	"shared/conf"
	"executionenvironment/executionenvironment"
	"os"
	"strconv"
	"shared/net"
	"shared/shared"
)

func main(){

	// Read OS arguments
	shared.ProcessOSArguments(os.Args[1:])

	// start configuration
	EE := executionenvironment.ExecutionEnvironment{}
	EE.Exec(conf.GenerateConf(parameters.DIR_CONF + "/MiddlewareQueueServer.conf"),parameters.IS_ADAPTIVE)

	fmt.Println("Queue server started at "+netshared.ResolveHostIp()+" Port= "+strconv.Itoa(parameters.QUEUESERVER_PORT))
	fmt.Scanln()
}
