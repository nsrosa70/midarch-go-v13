package main

import (
	"fmt"
	"shared/parameters"
	"strconv"
	"shared/net"
	EE "execution/execution"
)

func main(){

	// start configuration
	EE.ExecutionEnvironment{}.Deploy("QueueServer.conf")

	fmt.Println("Queue server started at "+netshared.ResolveHostIp()+" Port= "+strconv.Itoa(parameters.QUEUEING_PORT))
	fmt.Scanln()
}
