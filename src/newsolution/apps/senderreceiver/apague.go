package main

import (
	"newsolution/development/repositories/architectural"
	"reflect"
	"fmt"
	"newsolution/shared/shared"
)

type X struct{}

func (X) F(x int) {
	fmt.Println(x)
}

func main() {

	lib := architectural.ArchitecturalRepository{}
	lib.Load()

	obj := lib.Library["SRH"].Type

	msg := new(interface{})
	*msg = new(interface{})


	info := make([]*interface{}, 3)
	info[0] = new(interface{})
	*info[0] = new(interface{})
	info[1] = new(interface{})
	*info[1] = new(interface{})
	info[2] = new(interface{})
	*info[2] = new(interface{})

	fmt.Println(*info[0])
	shared.InvokeNewNew(obj, reflect.TypeOf(obj).Method(2).Name, msg, info) // I_Test1
	temp := *info[0]
	temp2 := temp.(int)
	fmt.Println(temp2)

	shared.InvokeNewNew(obj, reflect.TypeOf(obj).Method(3).Name, msg, info) // I_Test2
	temp = *info[0]
	temp2 = temp.(int)
	fmt.Println(temp2)
}

func oldmain() {

	lib := architectural.ArchitecturalRepository{}
	lib.Load()

	obj := lib.Library["SRH"].Type

	methods := []string{}

	for i := 0; i < reflect.TypeOf(obj).NumMethod(); i++ {
		methods = append(methods, reflect.TypeOf(obj).Method(i).Name)
	}
	fmt.Println(methods)

	for i := 0; i < len(methods); i++ {
		if methods[i] == "I_Test" {
			m := reflect.ValueOf(obj).MethodByName(methods[i])
			//in := make([]reflect.Value, m.Type().NumIn())
			in := make([]*reflect.Value, m.Type().NumIn())
			//in := make([]*interface{}, m.Type().NumIn())
			for i := 0; i < m.Type().NumIn(); i++ {
				t := m.Type().In(i)
				//in[i] = reflect.Value{}
				//in[i] = reflect.Zero(t)
				in[i] = new(reflect.Value)
				*in[i] = reflect.Zero(t)
			}
			fmt.Println(in[0])
			//reflect.ValueOf(obj).MethodByName("I_Test").Call(in)
			//shared.InvokeNew(obj,"I_Test",in)
			fmt.Println(in[0])
		}
	}
	//	t := lib.Library["SRH"].Type

	//	args := make([]interface{}, 3)
	//	args[0] = new(interface{})
	//	args[1] = new(interface{})
	//	args[2] = new(interface{})

	//msg := messages.SAMessage{}
	//args1 := make([]*interface{}, 1)
	//args1[0] = new(interface{})
	//*args1[0] = &msg

	//inputs := make([]reflect.Value, len(args1))
	//for i, _ := range args1 {
	//	inputs[i] = reflect.ValueOf(*args1[i])
	//}
	//reflect.ValueOf(t).MethodByName("I_Test").Call(inputs)
	//shared.InvokeOld(t, "I_Test", args1)

	//	fmt.Println(*args1[0])
	//fmt.Println(inputs[0])

	//	args2 := make([]reflect.Value,1)
	//x := reflect.Value{}
	//	args2[0] = reflect.Value{}
	//	args2[0] = reflect.Zero(tps[0])
	//args2[1] = reflect.Value{}
	//args2[1] = reflect.Zero(tps[1])
	//args2[2] = reflect.Value{}
	//args2[2] = reflect.Zero(tps[2])

	//fmt.Println(in[0])

	//fmt.Printf("Shared:: %v %v %v\n",reflect.TypeOf(any),name, inputs)

	//reflect.ValueOf(any).MethodByName(name).Call(inputs)

	//shared.InvokeNew(t,"I_Receive",args2)

	//	args[0] = &messages.SAMessage{Payload:"TODO"}
	//	args[1] = reflect.Zero(tps[1])
	//	args[2] = reflect.Zero(tps[2])
	//shared.Invoke(obj,methods[i],args)
	//shared.InvokeNew(obj,methods[i],args)
	//	os.Exit(0)
}

func initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			initializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
}
