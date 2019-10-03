package main

import (
	"newsolution/development/repositories/architectural"
	"reflect"
	"fmt"
	"gmidarch/development/framework/messages"
	"os"
	"newsolution/shared/shared"
)

type X struct{}

func (X) F(x int) {
	x = x + 1
}

func main() {

	lib := architectural.ArchitecturalRepository{}
	lib.Load()

	obj := lib.Library["SRH"].Type

	methods := []string{}

	for i := 0; i < reflect.TypeOf(obj).NumMethod(); i++{
		methods = append(methods,reflect.TypeOf(obj).Method(i).Name)
	}
	fmt.Println(methods)

	for i := 0; i < len(methods); i++ {
		m := reflect.ValueOf(obj).MethodByName(methods[i])
		in := make([]reflect.Value, m.Type().NumIn())
		tps := make([]reflect.Type,3)
		for i := 0; i < m.Type().NumIn(); i++ {
			t := m.Type().In(i)
			x := reflect.Zero(t)
			tps[i] = t
			fmt.Printf("HHH: %v\n",x.Type())
			in[i] = reflect.ValueOf(t)
		}


		t := lib.Library["SRH"].Type

		args := make([]interface{},3)
		args[0] = new(interface{})
		args[1] = new(interface{})
		args[2] = new(interface{})


		args2 := make([]reflect.Value,3)
		args2[0] = reflect.Zero(tps[0])
		args2[1] = reflect.Zero(tps[1])
		args2[2] = reflect.Zero(tps[2])

		fmt.Println(in[0])
		//reflect.ValueOf(t).MethodByName("I_Receive").Call(args2)
		shared.InvokeNew(t,"I_Receive",args2)


		args[0] = &messages.SAMessage{Payload:"TODO"}
		args[1] = reflect.Zero(tps[1])
		args[2] = reflect.Zero(tps[2])
		//shared.Invoke(obj,methods[i],args)
		//shared.InvokeNew(obj,methods[i],args)
		os.Exit(0)
	}
}
