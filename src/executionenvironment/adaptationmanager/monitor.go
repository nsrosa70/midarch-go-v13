package adaptationmanager

import (
	"time"
	"strings"
	"os"
	"io/ioutil"
	"shared/shared"
	"framework/configuration/configuration"
	"shared/parameters"
	"shared/errors"
)

type Monitor struct{}

func MonitorCorrective(chanInMonitoredCorrective chan shared.MonitoredCorrectiveData) {
	// TODO
}

func MonitorEvolutive(chanInMonitoredEvolutive chan shared.MonitoredEvolutiveData) {

	listOfOldPlugins := loadPlugins()
	for {
		listOfNewPlugins := loadPlugins()
		newPlugins := checkNewPlugins(listOfOldPlugins, listOfNewPlugins)
		if len(newPlugins) > 0 {
			chanInMonitoredEvolutive <- newPlugins
		}
		listOfOldPlugins = listOfNewPlugins
		time.Sleep(parameters.MONITOR_TIME * time.Second)
	}
}

func (Monitor) Exec(conf configuration.Configuration, chanMACorrective chan shared.MonitoredCorrectiveData, chanMAEvolutive chan shared.MonitoredEvolutiveData) {

	chanInMonitoredCorrective := make(chan shared.MonitoredCorrectiveData)
	chanInMonitoredEvolutive := make(chan shared.MonitoredEvolutiveData)

	//go MonitorCorrective(chanInMonitoredCorrective)
	go MonitorEvolutive(chanInMonitoredEvolutive)

	for {
		select {
		case monitoredData := <-chanInMonitoredCorrective:
			chanMACorrective <- monitoredData
		case listOfPlugins := <-chanInMonitoredEvolutive:
			chanMAEvolutive <- listOfPlugins
		}
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