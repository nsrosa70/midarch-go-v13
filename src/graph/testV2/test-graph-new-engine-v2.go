package main

import (
	"fmt"
	"reflect"
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
	//mainClientServer()
	//mainSenderReceiver()
	mainNamingServer()
}

func mainNamingServer() {
	conf := conf.GenerateConf("MiddlewareNamingServer.conf")
	fdrGraph := CreateFDRGraph()
	execGraph, channels := CreateExecGraph(fdrGraph)

	elemChannels1 := DefineChannels(channels, "srh")
	i_PreInvR1 := DefineChannel(elemChannels1, "I_PreInvR_srh")
	invR1 := DefineChannel(elemChannels1, "InvR")
	terR1 := DefineChannel(elemChannels1, "TerR")
	i_PosTerR1 := DefineChannel(elemChannels1, "I_PosTerR_srh")

	//elemChannels2 := DefineChannels(channels, "invoker")
	//invP2 := DefineChannel(elemChannels2, "InvP")
	//i_PosInvP2 := DefineChannel(elemChannels2, "I_PosInvP_invoker")
	//terP2 := DefineChannel(elemChannels2, "TerP")

	go Control(execGraph)
	go shared.Invoke(conf.Components["srh"].TypeElem, "Loop", i_PreInvR1,invR1,terR1,i_PosTerR1)
	//go shared.Invoke(conf.Components["invoker"].TypeElem, "Loop", invP2, i_PosInvP2,terP2)

	fmt.Scanln()
}


func CreateFDRGraph() fdrgraph.Graph {
	graph := fdrgraph.NewGraph(20)

	// MiddlewareNamingServer
	graph.AddEdge(0, 1, "I_PreInvR_srh")
	graph.AddEdge(1, 2, "InvR.srh")
	graph.AddEdge(2, 3, "InvP.invoker")
	graph.AddEdge(3, 4, "I_PosInvP_invoker")
	graph.AddEdge(4, 5, "TerP.invoker")
	graph.AddEdge(5, 6, "TerR.srh")
	graph.AddEdge(6, 0, "I_PosTerR_srh")

	// Sender/Receiver
	//graph.AddEdge(0, 1, "InvR.sender")
	//graph.AddEdge(1, 2, "InvP.receiver")
	//graph.AddEdge(2, 0, "I_PosInvP_receiver")
	//graph.AddEdge(2, 3, "InvR.sender")
	//graph.AddEdge(3, 1, "I_PosInvP_receiver")

	// Client/Server
	//graph.AddEdge(0, 1, "InvR.client")
	//graph.AddEdge(1, 2, "InvP.server")
	//graph.AddEdge(2, 3, "I_PosInvP_server")
	//graph.AddEdge(3, 4, "TerP.server")
	//graph.AddEdge(4, 0, "TerR.client")

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
	var value reflect.Value

	for i := 0; i < len(edges); i++ {
		if IsToElement(edges[i].Action.Action) {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(edges[i].Action.Channel), Send: reflect.ValueOf(*msg)}
		} else {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(edges[i].Action.Channel), Send: reflect.Value{}}
		}
	}

	*chosen, value, _ = reflect.Select(cases)
	if !IsToElement(edges[*chosen].Action.Action) {
		*msg = value.Interface().(string)
	}
	cases = nil
}

func IsToElement(action string) bool {
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
			node = edges[0].To
			if IsToElement(edges[0].Action.Action) {
				edges[0].Action.Channel <- msg
			} else {
				msg = <-edges[0].Action.Channel
			}
		} else { // two+ edges
			chosen := 0
			Choice(&msg, &chosen, edges)
			node = edges[chosen].To
		}
	}
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
	var r chan string
	found := false

	for c := range channels {
		if (a[:2] != "I_") {
			if strings.Contains(c, a) && c[:2] != "I_" {
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

func mainSenderReceiver() {
	conf := conf.GenerateConf("SenderReceiver.conf")
	fdrGraph := CreateFDRGraph()
	execGraph, channels := CreateExecGraph(fdrGraph)

	elemChannels1 := DefineChannels(channels, "sender")
	invR1 := DefineChannel(elemChannels1, "InvR")

	elemChannels2 := DefineChannels(channels, "receiver")
	invP2 := DefineChannel(elemChannels2, "InvP")
	i_PosInvP2 := DefineChannel(elemChannels2, "I_PosInvP_receiver")

	go Control(execGraph)
	go shared.Invoke(conf.Components["sender"].TypeElem, "Loop", invR1)
	go shared.Invoke(conf.Components["receiver"].TypeElem, "Loop", invP2, i_PosInvP2)

	fmt.Scanln()
}

func mainClientServer() {
	conf := conf.GenerateConf("ClientServer.conf")
	fdrGraph := CreateFDRGraph()
	execGraph, channels := CreateExecGraph(fdrGraph)

	elemChannels1 := DefineChannels(channels, "client")

	invR1 := DefineChannel(elemChannels1, "InvR")
	terR1 := DefineChannel(elemChannels1, "TerR")

	elemChannels2 := DefineChannels(channels, "server")
	invP2 := DefineChannel(elemChannels2, "InvP")
	i_PosInvP2 := DefineChannel(elemChannels2, "I_PosInvP_server")
	terP2 := DefineChannel(elemChannels2, "TerP")

	go Control(execGraph)
	go shared.Invoke(conf.Components["client"].TypeElem, "Loop", invR1, terR1)
	go shared.Invoke(conf.Components["server"].TypeElem, "Loop", invP2, terP2, i_PosInvP2)

	fmt.Scanln()
}
