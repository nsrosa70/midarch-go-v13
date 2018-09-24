package element

import (
	"framework/message"
	"graph/wgraph"
	"strings"
	"reflect"
	"fmt"
	"shared/shared"
)

// element definition
type Element struct {
	Id           string
	Behaviour    func(Element, wgraph.Graph, message.Message, map[string]chan message.Message, map[string]string)
	TypeElem     interface{}
	BehaviourExp string
}

// external actions common to all components and connectors
func (Element) InvP(invP *chan message.Message, msg *message.Message) {
	*msg = <-*invP
}

func (Element) InvR(invR *chan message.Message, msg *message.Message) {
	*invR <- *msg
}

func (Element) TerR(terR *chan message.Message, msg *message.Message) {
	*msg = <-*terR
}

func (Element) TerP(terP *chan message.Message, msg *message.Message) {
	*terP <- *msg
}

func choice(msg *message.Message, chosen *int, edges []wgraph.Edge) {
	cases := make([]reflect.SelectCase, len(edges))

	var value reflect.Value
	for i := 0; i < len(edges); i++ {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Action.P2)}
	}
	*chosen, value,_ = reflect.Select(cases)
	*msg = value.Interface().(message.Message)

	cases = nil
}

// behaviour execution
func (Element) BehaviourExec(elem Element, graph wgraph.Graph, msg message.Message, channs map[string]chan message.Message, elemMaps map[string]string) {

	// execute graph
	node := 0
	for {
		edges := graph.AdjacentEdges(node)
		if len(edges) == 1 {
			if shared.IsInternal(edges[0].Action.P4) {
				//Log(time.Now().String(),elem.Id,edges[0].Action.P4,"START") // TODO
				//edges[0].Action.F2(edges[0].Action.P3, edges[0].Action.P4, edges[0].Action.P1)
				edges[0].Action.F2(elem.TypeElem, edges[0].Action.P4, edges[0].Action.P1)
				//Log(time.Now().String(),elem.Id,edges[0].Action.P4,"complete") // TODO
			} else {
				//Log(time.Now().String(),elem.Id,edges[0].Action.P4,"START") // TODO
				edges[0].Action.F1(edges[0].Action.P2, edges[0].Action.P1)
				//Log(time.Now().String(),elem.Id,edges[0].Action.P4,"complete") // TODO
			}
			node = edges[0].To
		} else {
			msg := message.Message{}
			chosen := 0
			//Log(time.Now().String(),"CHOICE","START") // TODO
			choice(&msg, &chosen, edges)
			//Log(time.Now().String(),elem.Id,edges[chosen].Action.P4,"complete") // TODO
			*edges[chosen].Action.P1 = msg
			node = edges[chosen].To
		}
		if node == 0 {
			break
		}
	}
	return
}

func Log(args...string){
	if strings.Contains(args[1],"Proxy") || strings.Contains(args[1],"XXX"){
		fmt.Println(args[0]+":"+args[1]+":"+args[2]+":"+args[3])
	}
}