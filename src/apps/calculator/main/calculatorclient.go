package main

import (
	"framework/components/naming/naming"
	"fmt"
	"time"
	"shared/parameters"
	"executionenvironment"
	"apps/calculator/calculatorclientproxy"
	"shared/conf"
)

func main(){

	// start configuration
	EE := executionenvironment.ExecutionEnvironment{}
	EE.Exec(conf.GenerateConf(parameters.DIR_CONF + "/MiddlewareCalculatorClient.conf"),!parameters.IS_ADAPTIVE)

	// proxy to naming service
	namingClientProxy := naming.LocateNaming("localhost",parameters.NAMING_PORT)

	// obtain ior
	calc := namingClientProxy.Lookup("Calculator").(calculatorclientproxy.CalculatorClientProxy)

	// invoke remote method
	for i:= 0; i< parameters.SAMPLE_SIZE; i++ {

		t1 := time.Now()
		calc.Add(i,i)
		t2 := time.Now()

		x:= float64(t2.Sub(t1).Nanoseconds())/1000000

		fmt.Printf("%F \n",x)
	}

	/*
	// calculate meantime
	sd := float64(0)
	meanTime := float64(totalTime) / float64(shared.SAMPLE_SIZE)

	// standard deviation
	for i := 0; i < shared.SAMPLE_SIZE; i++ {
		// The use of Pow math function func Pow(x, y float64) float64
		sd += math.Pow(float64(t[i]) - meanTime, 2)
	}
	// The use of Sqrt math function func Sqrt(x float64) float64
	sd = math.Sqrt(sd/shared.SAMPLE_SIZE)

	fmt.Printf("Number of Invocations : %d \n",shared.SAMPLE_SIZE)
	fmt.Printf("Total time            : %v ms\n",totalTime/1000000)
	fmt.Printf("Average time          : %f ms\n",meanTime/1000000)
	fmt.Printf("Standard Deviation    : %f ms\n",sd/1000000)
	fmt.Printf("Service time          : %d ms \n", shared.SERVICE_TIME)
	fmt.Printf("Monitor time          : %d s \n", shared.MONITOR_TIME)
	fmt.Printf("Injection time        : %d s \n", shared.INJECTION_TIME)
	fmt.Printf("Injection Enabled     : %t \n", shared.INJECTION_ENABLED)
	fmt.Printf("Adaptive              : %t \n", shared.IS_ADAPTIVE)
	fmt.Printf("Strategy              : %d (2-Alternate 3-Non Alternate) \n", shared.STRATEGY)
	*/
}
