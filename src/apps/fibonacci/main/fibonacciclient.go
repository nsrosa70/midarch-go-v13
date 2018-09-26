package main

import (
	"framework/components/naming/naming"
	"fmt"
	"time"
	"apps/fibonacci/fibonacciclientproxy"
	"shared/parameters"
	EE "executionenvironment/executionenvironment"
)

func main(){

	// start configuration
	EE.ExecutionEnvironment{}.Exec("MiddlewareFibonacciClient.conf")

	// proxy to naming service
	namingClientProxy := naming.LocateNaming()

	// obtain ior
	fibo := namingClientProxy.Lookup("Fibonacci").(fibonacciclientproxy.FibonacciClientProxy)

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
