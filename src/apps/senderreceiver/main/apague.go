package main

import (
	"fmt"
	"reflect"
)

func main() {
	var channs = []chan int{
		make(chan int),
		make(chan int)}

	go choice(channs)

	fmt.Scanln()
}

func decide(channs [] chan int, msg int) {

	if msg == 1 {
		channs[0] <- 1
	} else {
		channs[1] <- 2
	}
}

func choice(channs [] chan int) {

	cases := make([]reflect.SelectCase, 2)

	var value reflect.Value
	for i := 0; i < 2; i++ {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channs[i])}
	}
	go decide(channs,2)
	_, value, _ = reflect.Select(cases)
	msg := value.Interface().(int)

	fmt.Println(msg)
	cases = nil
}
