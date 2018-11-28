package main

import (
	"fmt"
	"shared/parameters"
	"core/engine"
	"strconv"
	"shared/net"
)

func main() {

	// start configuration
	engine.ExecutionEnvironment{}.Deploy("MiddlewareNamingServer.conf") // TODO

	fmt.Println("Naming service started at " + netshared.ResolveHostIp() + " Port= " + strconv.Itoa(parameters.NAMING_PORT))
	fmt.Scanln()
}
