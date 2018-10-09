package executionenvironment

import (
	"shared/conf"
	"shared/shared"
	"os"
	"verificationtools/fdr"
	"framework/configuration/configuration"
	"framework/message"
	"graph/execgraph"
	"reflect"
	"framework/library"
	"strings"
)

type ExecutionEnvironment struct{}

func (ee ExecutionEnvironment) Deploy(confFile string) {

	// Load execution parameters
	shared.LoadParameters(os.Args[1:])

	// Generate Go configuration
	conf := conf.GenerateConf(confFile)
	fdrGraph := fdr.FDR{}.CreateFDRGraph(confFile)
	execGraph, channels := shared.CreateExecGraph(fdrGraph)
	ee.DeployGeneric(confFile,conf,channels,execGraph)
}

func (ee ExecutionEnvironment) DeployGeneric(confFile string, conf configuration.Configuration, channels map[string]chan message.Message, execGraph execgraph.GraphX) {

	go shared.Control(execGraph)
	for e := range conf.Components {
		elemChannels := shared.DefineChannels(channels, conf.Components[e].Id)
		actions := map[string][]string{}
		behaviour := library.Repository[reflect.TypeOf(conf.Components[e].TypeElem).String()].CSP
		actions[e] = FilterActions(strings.Split(behaviour," "))
		individualChannels := map[string]chan message.Message{}
		for a := range actions[e] {
			individualChannels[actions[e][a]] = shared.DefineChannel(elemChannels,actions[e][a])
		}
		go shared.Invoke(conf.Components[conf.Components[e].Id].TypeElem, "Loop", individualChannels)
	}
}

func (ee ExecutionEnvironment) DeployFibonacciServer(confFile string, conf configuration.Configuration, channels map[string]chan message.Message, execGraph execgraph.GraphX) {

	// naming proxy
	elemChannels1 := shared.DefineChannels(channels, "namingproxy")
	I_PreInvR_namingproxy1 := shared.DefineChannel(elemChannels1, "I_PreInvR_namingproxy")
	InvR1 := shared.DefineChannel(elemChannels1, "InvR.namingproxy")
	TerR1 := shared.DefineChannel(elemChannels1, "TerR")
	I_PosTerR_namingproxy1 := shared.DefineChannel(elemChannels1, "I_PosTerR_namingproxy")

	// invoker
	elemChannels2 := shared.DefineChannels(channels, "fibonacciinvoker")
	InvP2 := shared.DefineChannel(elemChannels2, "InvP")
	I_PosInvP2 := shared.DefineChannel(elemChannels2, "I_PosInvP_fibonacciinvoker")
	TerP2 := shared.DefineChannel(elemChannels2, "TerP")

	// requestor
	elemChannels3 := shared.DefineChannels(channels, "requestor")
	InvP3 := shared.DefineChannel(elemChannels3, "InvP")
	I_PosInvP_requestor3 := shared.DefineChannel(elemChannels3, "I_PosInvP_requestor")
	InvR3 := shared.DefineChannel(elemChannels3, "InvR")
	TerR3 := shared.DefineChannel(elemChannels3, "TerR")
	I_PosTerR_requestor3 := shared.DefineChannel(elemChannels3, "I_PosTerR_requestor")
	TerP3 := shared.DefineChannel(elemChannels3, "TerP")

	// crh
	elemChannels4 := shared.DefineChannels(channels, "crh")
	InvP4 := shared.DefineChannel(elemChannels4, "InvP")
	I_PosInvP_crh4 := shared.DefineChannel(elemChannels4, "I_PosInvP_crh")
	I_PreTerP_crh4 := shared.DefineChannel(elemChannels4, "I_PreTerP_crh")
	TerP4 := shared.DefineChannel(elemChannels4, "TerP")

	// srh
	elemChannels5 := shared.DefineChannels(channels, "srh")
	I_PreInvR5 := shared.DefineChannel(elemChannels5, "I_PreInvR_srh")
	InvR5 := shared.DefineChannel(elemChannels5, "InvR")
	TerR5 := shared.DefineChannel(elemChannels5, "TerR")
	I_PosTerR5 := shared.DefineChannel(elemChannels5, "I_PosTerR_srh")

	go shared.Control(execGraph)
	go shared.Invoke(conf.Components["namingproxy"].TypeElem, "Loop", I_PreInvR_namingproxy1, InvR1, TerR1, I_PosTerR_namingproxy1)
	go shared.Invoke(conf.Components["fibonacciinvoker"].TypeElem, "Loop", InvP2, I_PosInvP2, TerP2)
	go shared.Invoke(conf.Components["requestor"].TypeElem, "Loop", InvP3, I_PosInvP_requestor3, InvR3, TerR3, I_PosTerR_requestor3, TerP3)
	go shared.Invoke(conf.Components["crh"].TypeElem, "Loop", InvP4, I_PosInvP_crh4, I_PreTerP_crh4, TerP4)
	go shared.Invoke(conf.Components["srh"].TypeElem, "Loop", I_PreInvR5, InvR5, TerR5, I_PosTerR5)
}

func (ee ExecutionEnvironment) DeployFibonacciClient(confFile string, conf configuration.Configuration, channels map[string]chan message.Message, execGraph execgraph.GraphX) {

	// naming proxy
	elemChannels1 := shared.DefineChannels(channels, "namingproxy")
	I_PreInvR_namingproxy1 := shared.DefineChannel(elemChannels1, "I_PreInvR_namingproxy")
	InvR1 := shared.DefineChannel(elemChannels1, "InvR.namingproxy")
	TerR1 := shared.DefineChannel(elemChannels1, "TerR")
	I_PosTerR_namingproxy1 := shared.DefineChannel(elemChannels1, "I_PosTerR_namingproxy")

	// fibonacci proxy
	elemChannels2 := shared.DefineChannels(channels, "fibonacciproxy")
	I_PreInvR_fibonacciproxy2 := shared.DefineChannel(elemChannels2, "I_PreInvR_fibonacciproxy")
	InvR2 := shared.DefineChannel(elemChannels2, "InvR.fibonacciproxy")
	TerR2 := shared.DefineChannel(elemChannels2, "TerR")
	I_PosTerR_fibonacciproxy2 := shared.DefineChannel(elemChannels2, "I_PosTerR_fibonacciproxy")

	// requestor
	elemChannels3 := shared.DefineChannels(channels, "requestor")
	InvP3 := shared.DefineChannel(elemChannels3, "InvP")
	I_PosInvP_requestor3 := shared.DefineChannel(elemChannels3, "I_PosInvP_requestor")
	InvR3 := shared.DefineChannel(elemChannels3, "InvR")
	TerR3 := shared.DefineChannel(elemChannels3, "TerR")
	I_PosTerR_requestor3 := shared.DefineChannel(elemChannels3, "I_PosTerR_requestor")
	TerP3 := shared.DefineChannel(elemChannels3, "TerP")

	// crh
	elemChannels4 := shared.DefineChannels(channels, "crh")
	InvP4 := shared.DefineChannel(elemChannels4, "InvP")
	I_PosInvP_crh4 := shared.DefineChannel(elemChannels4, "I_PosInvP_crh")
	I_PreTerP_crh4 := shared.DefineChannel(elemChannels4, "I_PreTerP_crh")
	TerP4 := shared.DefineChannel(elemChannels4, "TerP")

	go shared.Control(execGraph)
	go shared.Invoke(conf.Components["namingproxy"].TypeElem, "Loop", I_PreInvR_namingproxy1, InvR1, TerR1, I_PosTerR_namingproxy1)
	go shared.Invoke(conf.Components["fibonacciproxy"].TypeElem, "Loop", I_PreInvR_fibonacciproxy2, InvR2, TerR2, I_PosTerR_fibonacciproxy2)
	go shared.Invoke(conf.Components["requestor"].TypeElem, "Loop", InvP3, I_PosInvP_requestor3, InvR3, TerR3, I_PosTerR_requestor3, TerP3)
	go shared.Invoke(conf.Components["crh"].TypeElem, "Loop", InvP4, I_PosInvP_crh4, I_PreTerP_crh4, TerP4)
}

func FilterActions(actions []string) [] string {
	r := []string{}

	for a := range actions {
		action := actions[a]
		if strings.Contains(action, "I") || strings.Contains(action,"T"){ // TODO
			if strings.Contains(action,"."){
				action = action[:strings.Index(action,".")]
			}
			r = append(r, action)
		}
	}
	return r
}