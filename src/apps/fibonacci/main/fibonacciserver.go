package main

import (
	"fmt"
	"strconv"
	"apps/fibonacci/fibonacciclientproxy"
	"shared/parameters"
	"executionenvironment/executionenvironment"
	"shared/net"
	"shared/factories"
)

func main(){

	// start configuration
	executionenvironment.ExecutionEnvironment{}.Deploy("MiddlewareFibonacciServer.conf")

	// proxy to naming service
	namingClientProxy := factories.LocateNaming()

	// register
	fiboProxy := fibonacciclientproxy.FibonacciClientProxy{Host:netshared.ResolveHostIp(),Port:parameters.FIBONACCI_PORT} // TODO
	namingClientProxy.Register("Fibonacci", fiboProxy)

	fmt.Println("Fibonacci Server ready at port "+strconv.Itoa(fiboProxy.Port))

	fmt.Scanln()
	fmt.Println("done")
}
