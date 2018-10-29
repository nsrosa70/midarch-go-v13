package proxy

import (
	"reflect"
	"framework/components"
)

var ProxyLibrary = map[string]reflect.Type{
	"components.CalculatorClientProxy": reflect.TypeOf(components.CalculatorClientProxy{}),
	"components.FibonacciClientProxy":  reflect.TypeOf(components.FibonacciClientProxy{})}
