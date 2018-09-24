package main

import (
	"fmt"
	"strconv"
	"framework/components/naming/naming"
	"executionenvironment"
	"shared/parameters"
	"apps/calculator/calculatorclientproxy"
	"shared/conf"
)

func main(){
	// proxy to naming service
	namingClientProxy := naming.LocateNaming("localhost",parameters.NAMING_PORT)

	// start configuration
	EE := executionenvironment.ExecutionEnvironment{}
	EE.Exec(conf.GenerateConf(parameters.DIR_CONF + "/MiddlewareCalculatorServer.conf"),parameters.IS_ADAPTIVE)

	// register
	calculator := calculatorclientproxy.CalculatorClientProxy{Host:"localhost",Port:parameters.CALCULATOR_PORT}
	namingClientProxy.Register("Calculator", calculator)

	fmt.Println("Calculator Server ready at port "+strconv.Itoa(calculator.Port))

	fmt.Scanln()
	fmt.Println("done")
}
