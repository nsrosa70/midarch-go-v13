package main

import (
	"fmt"
	"strconv"
	"executionenvironment/executionenvironment"
	"shared/net"
	"shared/factories"
	"shared/parameters"
	"framework/components"
)

func main(){

	// Parameters: NAMING_HOST, STRATEGY, MONITOR_TIME, INJECTION_TIME
	// start configuration
	executionenvironment.ExecutionEnvironment{}.Deploy("MiddlewareFibonacciServer.conf")

	// proxy to naming service
	namingClientProxy := factories.LocateNaming()

	// register
	fiboProxy := components.FibonacciClientProxy{Host:netshared.ResolveHostIp(),Port:parameters.FIBONACCI_PORT}
	namingClientProxy.Register("Fibonacci", fiboProxy)

	fmt.Println("Fibonacci Server ready at port "+strconv.Itoa(fiboProxy.Port))

	fmt.Scanln()
	fmt.Println("done")
}
