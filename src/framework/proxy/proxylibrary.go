package proxy

import (
	"reflect"
	"framework/components"
)

var ProxyLibrary = map[string] reflect.Type {
	"calculatorclientproxy.CalculatorClientProxy": reflect.TypeOf(components.CalculatorClientProxy{}),
	"fibonacciclientproxy.FibonacciClientProxy": reflect.TypeOf(components.FibonacciClientProxy{})}
