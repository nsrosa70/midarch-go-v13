package executionenvironment

import (
	"shared/conf"
	"shared/shared"
	"os"
	"fmt"
	"verificationtools/fdr"
)

type ExecutionEnvironment struct{}

func (ee ExecutionEnvironment) Deploy(confFile string) {

	// Load execution parameters
	shared.LoadParameters(os.Args[1:])

	// Generate Go configuration
	conf := conf.GenerateConf("MiddlewareNamingServer.conf")
	fdrGraph := fdr.FDR{}.CreateFDRGraph()
	execGraph, channels := shared.CreateExecGraph(fdrGraph)


	// srh
	elemChannels1 := shared.DefineChannels(channels, "srh")
	i_PreInvR1 := shared.DefineChannel(elemChannels1, "I_PreInvR_srh")
	invR1 := shared.DefineChannel(elemChannels1, "InvR")
	terR1 := shared.DefineChannel(elemChannels1, "TerR")
	i_PosTerR1 := shared.DefineChannel(elemChannels1, "I_PosTerR_srh")

	// invoker
	elemChannels2 := shared.DefineChannels(channels, "invoker")
	invP2 := shared.DefineChannel(elemChannels2, "InvP")
	i_PosInvP2 := shared.DefineChannel(elemChannels2, "I_PosInvP_invoker")
	terP2 := shared.DefineChannel(elemChannels2, "TerP")

	go shared.Control(execGraph)
	go shared.Invoke(conf.Components["srh"].TypeElem, "Loop", i_PreInvR1,invR1,terR1,i_PosTerR1)
	go shared.Invoke(conf.Components["invoker"].TypeElem, "Loop", invP2, i_PosInvP2,terP2)

	fmt.Scanln()
	/*
	// Initialize channels between Units and Adaptation manager
	channsUnits := make(map[string]chan commands.LowLevelCommand)
	for i := range conf.Components {
		id := conf.Components[i].Id
		channsUnits[id] = make(chan commands.LowLevelCommand)
	}
	for i := range conf.Connectors {
		id := conf.Connectors[i].Id
		channsUnits[id] = make(chan commands.LowLevelCommand)
	}

	// Initialize basic elements used throughout execution
	channs := new(map[string]chan message.Message)
	elemMaps := new(map[string]string)

	// Configure channels & maps
	*channs = ee.ConfigureChannels(conf)
	*elemMaps = ee.ConfigureMaps(conf)

	// Start execution engine (manage the execution units)
	executionEngine := executionengine.ExecutionEngine{}
	go executionEngine.Exec(conf, *channs, *elemMaps, channsUnits)

	// Start adaptation manager
	if parameters.IS_ADAPTIVE {
		adaptationManager := adaptationmanager.AdaptationManager{}
		go adaptationManager.Exec(conf, *channs, *elemMaps, channsUnits)
		go versioninginjector.InjectAdaptiveEvolution(parameters.PLUGIN_BASE_NAME)
	}
	*/
}

