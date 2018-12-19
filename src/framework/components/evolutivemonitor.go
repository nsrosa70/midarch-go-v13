package components

import (
"framework/messages"
"time"
"shared/parameters"
"io/ioutil"
"shared/errors"
"strings"
"os"
"shared/shared"
)

type EvolutiveMonitor struct{}

var listOfOldPlugins = loadPlugins()

func (EvolutiveMonitor) I_EvolutiveMonitoring(msg *messages.SAMessage, info interface{}, r *bool) {

	//listOfOldPlugins := loadPlugins()

	listOfNewPlugins := loadPlugins()
	newPlugins := checkNewPlugins(listOfOldPlugins, listOfNewPlugins)
	if len(newPlugins) > 0 {
		evolutiveMonitoredData := shared.MonitoredEvolutiveData{}
		evolutiveMonitoredData = newPlugins
		*msg = messages.SAMessage{evolutiveMonitoredData}
	}
	listOfOldPlugins = listOfNewPlugins
	time.Sleep(parameters.MONITOR_TIME * time.Second)
}

func (EvolutiveMonitor) I_HasPlugin (msg *messages.SAMessage, info interface{}, r *bool) {

	if msg.Payload != nil{
		*r = true
	}
}

func (EvolutiveMonitor) I_HasNotPlugin (msg *messages.SAMessage, info interface{}, r *bool) {

	if msg.Payload == nil{
		*r = true
	}
}

func loadPlugins() map[string]time.Time {
	listOfPlugins := make(map[string]time.Time)

	pluginsDir, err := ioutil.ReadDir(parameters.DIR_PLUGINS)
	if err != nil {
		myError := errors.MyError{Source: "Analyser", Message: "Folder not read"}
		myError.ERROR()
	}
	for i := range pluginsDir {
		fileName := pluginsDir[i].Name()
		if strings.Contains(fileName, "_plugin") {
			info, err := os.Stat(parameters.DIR_PLUGINS + "/" + fileName)
			if err != nil {
				myError := errors.MyError{Source: "Analyser", Message: "Plugins not read"}
				myError.ERROR()
			}
			listOfPlugins[fileName] = info.ModTime()
		}
	}
	return listOfPlugins
}

func checkNewPlugins(listOfOldPlugins map[string]time.Time, listOfNewPlugins map[string]time.Time) [] string {
	var newPlugins [] string

	// check new plugins
	for key := range listOfNewPlugins {
		val1, _ := listOfNewPlugins[key]
		val2, ok2 := listOfOldPlugins[key]
		if ok2 {
			if val1.After(val2) { // newer version of an old plugin is available
				newPlugins = append(newPlugins, key)
			}
		} else {
			newPlugins = append(newPlugins, key) // a new plugin is available
		}
	}
	return newPlugins
}
