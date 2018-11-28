package main

import (
	"fmt"
	"time"
	"shared/parameters"
	EE "core/engine"
	"shared/factories"
	"framework/components"
)

func main(){

	// HOST_NAMING, SAMPLE_SIZE
	// start configuration
	EE.ExecutionEnvironment{}.Deploy("MiddlewareFibonacciClient.conf")

	// proxy to naming service
	namingClientProxy := factories.LocateNaming()

	// obtain proxy
	fibo := namingClientProxy.Lookup("Fibonacci").(components.FibonacciClientProxy)

	// invoke remote method
	for i:= 0; i< parameters.SAMPLE_SIZE; i++ {

		t1 := time.Now()
		fibo.Fibo(38)
		t2 := time.Now()

		x:= float64(t2.Sub(t1).Nanoseconds())/1000000
		fmt.Printf("%F \n",x)
		time.Sleep(parameters.REQUEST_TIME * time.Millisecond)
	}
}
