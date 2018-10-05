package main

import (
	"fmt"
	"graph/wgraph"
	"reflect"
	"strconv"
	"strings"
	"shared/shared"
)

type Sender struct{}
type Receiver struct{}
type Client struct{}
type Server struct{}

func main() {
	invR := make(chan string)
	invP := make(chan string, )
	terP := make(chan string)
	terR := make(chan string)
	i_PosInvP := make(chan string)

	graph := createGraph(&invR, &terR, &invP, &terP, &i_PosInvP)
	//go Client{}.Loop(invR, terR)
	//go Server{}.Loop(invP, terP, i_PosInvP)
	go Receiver{}.Loop(invP, i_PosInvP)
	go Sender{}.Loop(invR)
	go Control(graph)

	fmt.Scanln()
}

func createGraph(invR, terR, invP, terP, i_PosInvP *chan string) wgraph.Graph {
	graph := wgraph.NewGraph(20)

	// sender/receiver
	graph.AddEdge(0, 1, shared.ParMapActions{P5: invR, P4: "InvR.Sender"})
	graph.AddEdge(1, 2, shared.ParMapActions{P5: invP, P4: "InvP.Receiver"})
	graph.AddEdge(2, 0, shared.ParMapActions{P5: i_PosInvP, P4: "I_PosInvP_Receiver"})
	graph.AddEdge(2, 3, shared.ParMapActions{P5: invR, P4: "InvR.Sender"})
	graph.AddEdge(3, 1, shared.ParMapActions{P5: i_PosInvP, P4: "I_PosInvP_Receiver"})

	// clientX/server
	/*
	graph.AddEdge(0, 1, shared.ParMapActions{P5: invR, P4: "InvR.Client"})
	graph.AddEdge(1, 2, shared.ParMapActions{P5: invP, P4: "InvP.Server"})
	graph.AddEdge(2, 3, shared.ParMapActions{P5: i_PosInvP, P4: "I_PosInvP_Server"})
	graph.AddEdge(3, 4, shared.ParMapActions{P5: terP, P4: "TerP.Server"})
	graph.AddEdge(4, 0, shared.ParMapActions{P5: terR, P4: "TerR.Client"})
*/
	return *graph
}

func Choice(msg *string, chosen *int, edges []wgraph.Edge) {
	cases := make([]reflect.SelectCase, len(edges))

	for i := 0; i < len(edges); i++ {
		if IsSendAction(edges[i].Action.P4) {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(*edges[i].Action.P5), Send: reflect.ValueOf(*msg)}
		} else {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Action.P5), Send: reflect.Value{}}
		}
	}

	var value reflect.Value
	*chosen, value, _ = reflect.Select(cases)
	if !IsSendAction(edges[*chosen].Action.P4) {
		*msg = value.Interface().(string)
	}
	cases = nil
}

func IsSendAction(action string) bool {
	if action[:2] == "I_" || action[:4] == "InvP" || action[:4] == "TerR" {
		return true
	} else { // TerP and InvR
		return false
	}
}

func Control(g wgraph.Graph) {
	node := 0
	var msg = ""
	for {
		edges := g.AdjacentEdges(node)
		if len(edges) == 1 { // one edge
			chn := make(chan string)
			chn = *edges[0].Action.P5
			node = edges[0].To
			action := edges[0].Action.P4
			if !IsSendAction(action) {
				msg = <-chn
			} else {
				chn <- msg
			}
		} else { // two+ edges
			chosen := 0
			Choice(&msg, &chosen, edges)
			node = edges[chosen].To
		}
	}
}

func (Receiver) I_PosInvP(m *string) {
	fmt.Println("Receiver:: " + *m)
}

func (Receiver) Loop(invP, i_PosInvP chan string) {
	for {
		select {
		case <-invP:
		case msg := <-i_PosInvP:
			Receiver{}.I_PosInvP(&msg)
		}
	}
}

func (Sender) Loop(invR chan string) {
	msg := "testV1"
	i := 0
	for {
		invR <- msg + strconv.Itoa(i)
		i++
	}
}

func (Client) Loop(invR, terR chan string) {
	msgSent := "testV1"
	i := 0
	for {
		select {
		case invR <- msgSent + strconv.Itoa(i):
			i++
		case msgRecv := <-terR:
			fmt.Println("Message: " + msgRecv)
		}
	}
}

func (Server) Loop(invP, terP, i_PosInvP chan string) {
	msgRecv := ""
	for {
		select {
		case <-invP:
		case terP <- msgRecv:
		case msgRecv = <-i_PosInvP:
			Server{}.I_PosInvP(&msgRecv)
		}
	}
}

func (Server) I_PosInvP(msg *string) {
	*msg = strings.ToUpper(*msg)
}


