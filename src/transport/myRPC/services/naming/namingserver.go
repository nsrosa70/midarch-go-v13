package main

import (
	"fmt"
	"shared/parameters"
	"ee/ee"
	"strconv"
	"shared/net"
)

func main() {

	// start configuration
	ee.Engine{}.Deploy("MiddlewareNamingServer.conf") // TODO

	fmt.Println("Naming service started at " + netshared.ResolveHostIp() + " Port= " + strconv.Itoa(parameters.NAMING_PORT))
	fmt.Scanln()
}
