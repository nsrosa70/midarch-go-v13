package main

import (
	"fmt"
	"shared/parameters"
	"executionenvironment/executionenvironment"
	"strconv"
	"shared/net"
)

func main() {

	// start configuration
	executionenvironment.ExecutionEnvironment{}.Exec("MiddlewareNamingServer.conf")

	fmt.Println("Naming service started at " + netshared.ResolveHostIp() + " Port= " + strconv.Itoa(parameters.NAMING_PORT))
	fmt.Scanln()
}
