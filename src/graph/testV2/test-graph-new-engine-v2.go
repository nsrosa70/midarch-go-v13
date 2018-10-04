package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"graph/fdrgraph"
	"graph/execgraph"
	"shared/conf"
	"shared/shared"
)

type Receiver struct{}
type Sender struct{}
type Client struct{}
type Server struct{}

func main() {
	conf := conf.GenerateConf("SenderReceiver.conf")
	fdrGraph := CreateFDRGraph()
	execGraph, channels := CreateExecGraph(fdrGraph)

	go Control(execGraph)

	fmt.Println(channels)
	elemChannels1 := DefineChannels(channels, "sender")
	invR1 := DefineChannel(elemChannels1, "InvR")

	elemChannels2 := DefineChannels(channels, "receiver")
	invP := DefineChannel(elemChannels2, "InvP")
	i_PosInvP := DefineChannel(elemChannels2, "I_PosInvP")

	go shared.Invoke(conf.Components["sender"].TypeElem, "Loop", invR1)
	go shared.Invoke(conf.Components["receiver"].TypeElem, "Loop", invP, i_PosInvP)

	fmt.Scanln()
}

func DefineChannels(channels map[string]chan string, elem string) map[string]chan string {
	r := map[string]chan string{}

	for c := range channels {
		if strings.Contains(c, elem) {
			r[c] = channels[c]
		}
	}
	return r
}

func DefineChannel(channels map[string]chan string, a string) chan string {
	r := make(chan string)
	found := false

	for c := range channels {
		if (a[:2] != "I_") {
			if strings.Contains(c, a) {
				r = channels[c]
				found = true
				break
			}
		} else {
			if strings.Contains(c, a) {
				r = channels[c]
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println("Error: channel '" + a + " not found")
	}

	return r
}
func CreateFDRGraph() fdrgraph.Graph {
	graph := fdrgraph.NewGraph(20)

	// Sender/Receiver
	graph.AddEdge(0, 1, "InvR.sender")
	graph.AddEdge(1, 2, "InvP.receiver")
	graph.AddEdge(2, 0, "I_PosInvP_receiver")
	graph.AddEdge(2, 3, "InvR.sender")
	graph.AddEdge(3, 1, "I_PosInvP_receiver")

	// Client/Server
	graph.AddEdge(0, 1, "InvR.sender")
	graph.AddEdge(1, 2, "InvP.receiver")
	graph.AddEdge(2, 0, "I_PosInvP_receiver")
	graph.AddEdge(2, 3, "InvR.sender")
	graph.AddEdge(3, 1, "I_PosInvP_receiver")

	return *graph
}

func CreateExecGraph(fdrGraph fdrgraph.Graph) (execgraph.Graph, map[string]chan string) {
	graph := execgraph.NewGraph(fdrGraph.NumNodes)
	channels := map[string]chan string{}

	// create channels
	for e1 := range fdrGraph.Edges {
		for e2 := range fdrGraph.Edges[e1] {
			eTemp := fdrGraph.Edges[e1][e2]
			if _, ok := channels[eTemp.Action]; !ok {
				channels[eTemp.Action] = make(chan string)
			}
			graph.AddEdge(eTemp.From, eTemp.To, execgraph.ExecAction{Action: eTemp.Action, Channel: channels[eTemp.Action]})
		}
	}
	return *graph, channels
}

func Choice(msg *string, chosen *int, edges []execgraph.Edge) {
	cases := make([]reflect.SelectCase, len(edges))

	for i := 0; i < len(edges); i++ {
		if IsSendAction(edges[i].Action.Action) {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(edges[i].Action.Channel), Send: reflect.ValueOf(*msg)}
		} else {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(edges[i].Action.Channel), Send: reflect.Value{}}
		}
	}

	var value reflect.Value
	*chosen, value, _ = reflect.Select(cases)
	if !IsSendAction(edges[*chosen].Action.Action) {
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

func Control(g execgraph.Graph) {
	node := 0
	var msg = ""
	for {
		edges := g.AdjacentEdges(node)
		if len(edges) == 1 { // one edge
			chn := make(chan string)
			chn = edges[0].Action.Channel
			node = edges[0].To
			action := edges[0].Action.Action
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
