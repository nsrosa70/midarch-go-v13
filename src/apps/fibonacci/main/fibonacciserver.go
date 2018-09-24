package main

import (
	"fmt"
	"strconv"
	"framework/components/naming/naming"
	"apps/fibonacci/fibonacciclientproxy"
	"shared/parameters"
	"os"
	"shared/shared"
	"shared/conf"
	"executionenvironment/executionenvironment"
	"shared/net"
)

func main(){

	// read OS arguments
	shared.ProcessOSArguments(os.Args[1:])

	// start configuration
	EE := executionenvironment.ExecutionEnvironment{}
	EE.Exec(conf.GenerateConf(parameters.DIR_CONF + "/MiddlewareFibonacciServer.conf"),parameters.IS_ADAPTIVE)

	// proxy to naming service
	namingClientProxy := naming.LocateNaming(parameters.NAMING_HOST,parameters.NAMING_PORT)

	// register
	fibo := fibonacciclientproxy.FibonacciClientProxy{Host:netshared.ResolveHostIp(),Port:parameters.FIBONACCI_PORT} // TODO
	namingClientProxy.Register("Fibonacci", fibo)
	fmt.Println("Fibonacci Server ready at port "+strconv.Itoa(fibo.Port))

	fmt.Scanln()
	fmt.Println("done")
}
