package proxy

import (
	"reflect"
	"apps/calculator/calculatorclientproxy"
	"apps/fibonacci/fibonacciclientproxy"
)

var ProxyLibrary = map[string] reflect.Type {
	"calculatorclientproxy.CalculatorClientProxy": reflect.TypeOf(calculatorclientproxy.CalculatorClientProxy{}),
	"fibonacciclientproxy.FibonacciClientProxy": reflect.TypeOf(fibonacciclientproxy.FibonacciClientProxy{})}
