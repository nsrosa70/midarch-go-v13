package main

import (
	"fmt"
	"strconv"
	"core/engine"
	"shared/net"
	"shared/factories"
	"shared/parameters"
	"framework/components"
)

func main(){

	// Parameters: NAMING_HOST, STRATEGY, MONITOR_TIME, INJECTION_TIME
	// start configuration
	engine.Engine{}.Deploy("MiddlewareFibonacciServer.conf")

	// proxy to naming service
	namingClientProxy := factories.LocateNaming()

	// register
	fiboProxy := components.FibonacciClientProxy{Host:netshared.ResolveHostIp(),Port:parameters.FIBONACCI_PORT}
	namingClientProxy.Register("Fibonacci", fiboProxy)

	fmt.Println("Fibonacci Server ready at port "+strconv.Itoa(fiboProxy.Port))

	fmt.Scanln()
	fmt.Println("done")
}
