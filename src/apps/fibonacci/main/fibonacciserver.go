package main

import (
	"fmt"
	"strconv"
	"framework/components/naming/naming"
	"apps/fibonacci/fibonacciclientproxy"
	"shared/parameters"
	"executionenvironment/executionenvironment"
	"shared/net"
)

func main(){

	// start configuration
	executionenvironment.ExecutionEnvironment{}.Exec("MiddlewareFibonacciServer.conf")

	// proxy to naming service
	namingClientProxy := naming.LocateNaming(parameters.NAMING_HOST)

	// register
	fibo := fibonacciclientproxy.FibonacciClientProxy{Host:netshared.ResolveHostIp(),Port:parameters.FIBONACCI_PORT} // TODO
	namingClientProxy.Register("Fibonacci", fibo)
	fmt.Println("Fibonacci Server ready at port "+strconv.Itoa(fibo.Port))

	fmt.Scanln()
	fmt.Println("done")
}
