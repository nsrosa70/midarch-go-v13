package main

import (
	"shared/shared"
	"os"
	"executionenvironment/executionenvironment"
	"shared/conf"
	"shared/parameters"
	"framework/components/queueing/queueing"
)

func main() {

	// read OS arguments
	shared.ProcessOSArguments(os.Args[1:])

	// start configuration
	EE := executionenvironment.ExecutionEnvironment{}
	EE.Exec(conf.GenerateConf(parameters.DIR_CONF+"/"+"MiddlewarePublisher.conf"), parameters.IS_ADAPTIVE)

	// proxy to naming service
	queueingClientProxy := queueing.LocateQueueing(parameters.QUEUEING_HOST, parameters.QUEUESERVER_PORT)

	queueingClientProxy.Publish("topic1","teste")

	// obtain ior
	//fibo := namingClientProxy.Lookup("Fibonacci").(fibonacciclientproxy.FibonacciClientProxy)

	// invoke remote method
	//for i:= 0; i< parameters.SAMPLE_SIZE; i++ {

	//	t1 := time.Now()
	//	fibo.Fibo(38)
	//	t2 := time.Now()

	//	x:= float64(t2.Sub(t1).Nanoseconds())/1000000

	//	fmt.Printf("%F \n",x)

	//	time.Sleep(parameters.REQUEST_TIME * time.Millisecond)
	//}
}
