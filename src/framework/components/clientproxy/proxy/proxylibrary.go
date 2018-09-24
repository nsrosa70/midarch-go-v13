package proxy

import (
	"reflect"
	"apps/calculator/calculatorclientproxy"
	"apps/fibonacci/fibonacciclientproxy"
)

var ProxyLibrary = map[string] reflect.Type {
	reflect.TypeOf(calculatorclientproxy.CalculatorClientProxy{}).String(): reflect.TypeOf(calculatorclientproxy.CalculatorClientProxy{}),
	reflect.TypeOf(fibonacciclientproxy.FibonacciClientProxy{}).String(): reflect.TypeOf(fibonacciclientproxy.FibonacciClientProxy{})}

