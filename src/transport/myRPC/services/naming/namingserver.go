package main

import (
	"fmt"
	"gmidarch/shared/parameters"
	"gmidarch/execution/ee"
	"strconv"
	"gmidarch/shared/net"
)

func main() {

	// start configuration
	ee.Engine{}.Deploy("MiddlewareNamingServer.conf") // TODO

	fmt.Println("Naming service started at " + netshared.ResolveHostIp() + " Port= " + strconv.Itoa(parameters.NAMING_PORT))
	fmt.Scanln()
}
