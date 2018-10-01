package executionenvironment

import (
	"executionenvironment/executionengine"
	"framework/configuration/configuration"
	"framework/configuration/commands"
	"executionenvironment/adaptationmanager"
	"shared/parameters"
	"framework/message"
	"strings"
	"strconv"
	"executionenvironment/versioninginjector"
	"shared/conf"
	"shared/shared"
	"os"
)

type ExecutionEnvironment struct{}

func (ee ExecutionEnvironment) Deploy(confFile string) {

	// Load execution parameters
	shared.LoadParameters(os.Args[1:])

	// Generate Go configuration
	conf := conf.GenerateConf(confFile)

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

	// Start execution engine
	executionEngine := executionengine.ExecutionEngine{}
	go executionEngine.Exec(conf, *channs, *elemMaps, channsUnits)

	// Start adaptation manager
	if parameters.IS_ADAPTIVE {
		adaptationManager := adaptationmanager.AdaptationManager{}
		go adaptationManager.Exec(conf, *channs, *elemMaps, channsUnits)
		go versioninginjector.InjectAdaptiveEvolution(parameters.PLUGIN_BASE_NAME)
	}
}

func (ExecutionEnvironment) ConfigureChannels(conf configuration.Configuration) (map[string]chan message.Message) {
	channs := make(map[string]chan message.Message)

	for i := range conf.Attachments {
		c1Id := conf.Attachments[i].C1.Id
		c2Id := conf.Attachments[i].C2.Id
		tId := conf.Attachments[i].T.Id

		// c1 -> t
		key01 := c1Id + "." + "InvR" + "." + tId
		key02 := tId + "." + "InvP" + "." + c1Id
		key03 := tId + "." + "TerP" + "." + c1Id
		key04 := c1Id + "." + "TerR" + "." + tId
		channs[key01] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
		channs[key02] = channs[key01]
		channs[key03] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
		channs[key04] = channs[key03]

		// t -> c2
		key01 = tId + "." + "InvR" + "." + c2Id
		key02 = c2Id + "." + "InvP" + "." + tId
		key03 = c2Id + "." + "TerP" + "." + tId
		key04 = tId + "." + "TerR" + "." + c2Id
		channs[key01] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
		channs[key02] = channs[key01]
		channs[key03] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
		channs[key04] = channs[key03]
	}

	return channs
}

func (ExecutionEnvironment) ConfigureMaps(conf configuration.Configuration) (map[string]string) {

	elemMaps := make(map[string]string)
	partners := make(map[string]string)

	for i := range conf.Attachments {
		c1Id := conf.Attachments[i].C1.Id
		c2Id := conf.Attachments[i].C2.Id
		tId := conf.Attachments[i].T.Id
		if !strings.Contains(partners[c1Id], tId) {
			partners[c1Id] += ":" + tId
		}
		if !strings.Contains(partners[tId], c1Id) {
			partners[tId] += ":" + c1Id
		}
		if !strings.Contains(partners[tId], c2Id) {
			partners[tId] += ":" + c2Id
		}
		if !strings.Contains(partners[c2Id], tId) {
			partners[c2Id] += ":" + tId
		}
	}

	for i := range partners {
		p := strings.Split(partners[i], ":")
		c := 1
		for j := range p {
			if p[j] != "" {
				elemMaps[i+".e"+strconv.Itoa(c)] = p[j]
				c++
			}
		}
	}
	return elemMaps
}
