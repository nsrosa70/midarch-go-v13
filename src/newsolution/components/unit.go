package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/messages"
	"fmt"
	"reflect"
	"gmidarch/shared/shared"
)

type Unit struct {
	CSP string
	Graph graphs.GraphExecutable
	Buffer messages.SAMessage
	Element interface{}
}

func NewUnit(elem interface{}) Unit{
	r := new(Unit)

	r.CSP = "B = I_Init -> B"
	r.Buffer = messages.SAMessage{Payload:""}

	r.Graph = *graphs.NewGraph(1)
	newEdgeInfo := graphs.EdgeExecutableInfo{ActionName:"I_Execute",InternalAction:shared.Invoke}
	r.Graph.AddEdge(0,0,newEdgeInfo)

	return *r
}

func (u Unit) I_Execute(){
	//e := engine.Engine{}

	fmt.Println(reflect.TypeOf(u.Element))
	//e.Execute(u.Element)
}