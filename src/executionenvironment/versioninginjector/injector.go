package versioninginjector

import (
	"framework/library"
	"strings"
	"fmt"
	"os"
	"os/exec"
	"shared/parameters"
	"time"
)

func InjectAdaptiveEvolution(elementName string) {

	pluginBase01 := elementName + "01"
	sourceCode01 := pluginBase01 + ".go"

	pluginBase02 := elementName + "02"
	sourceCode02 := pluginBase02 + ".go"

	// remove old plugins
	outputLS, err := exec.Command("/bin/ls", parameters.DIR_PLUGINS).CombinedOutput()
	if err != nil {
		fmt.Println("Shared:: Something wrong in dir '"+parameters.DIR_PLUGINS)
		os.Stderr.WriteString(err.Error())
	}
	oldPlugins := strings.Split(string(outputLS), "\n")

	for plugin := range oldPlugins {
		if strings.Contains(oldPlugins[plugin], "_plugin_") {
			exec.Command("/bin/rm", "-r", parameters.DIR_PLUGINS+"/"+strings.TrimSpace(oldPlugins[plugin])).CombinedOutput()
		}
	}

	// generate new plugin
	switch parameters.STRATEGY {
	case 1: // no change
	case 2: // change once
		pluginName := strings.TrimSpace(pluginBase01 + "_plugin_v1")
		_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", parameters.DIR_PLUGINS+"/"+pluginName, parameters.DIR_PLUGINS+"/"+pluginBase01+"/"+sourceCode01).CombinedOutput()
		if err != nil {
			fmt.Println("Shared:: Something wrong in generating plugin '"+pluginName+"'")
			os.Stderr.WriteString(err.Error())
		}
	case 3: // change same plugin
		for {
			pluginName := strings.TrimSpace(pluginBase01 + "_plugin_v1")
			_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", parameters.DIR_PLUGINS+"/"+pluginName, parameters.DIR_PLUGINS+"/"+pluginBase01+"/"+sourceCode01).CombinedOutput()
			if err != nil {
				fmt.Println("Shared:: Something wrong in generating plugin '"+pluginName+"'")
				os.Stderr.WriteString(err.Error())
			}
			time.Sleep(parameters.INJECTION_TIME * time.Second)
		}
	case 4: // alternate plugins
		currentPlugin := 1
		for {
			switch currentPlugin {
			case 1:
				currentPlugin = 2
				pluginName := strings.TrimSpace(pluginBase01 + "_plugin_v1")
				_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", parameters.DIR_PLUGINS+"/"+pluginName, parameters.DIR_PLUGINS+"/"+pluginBase01+"/"+sourceCode01).CombinedOutput()
				if err != nil {
					fmt.Println("Shared:: Something wrong in generating plugin '"+pluginName+"'")
					os.Stderr.WriteString(err.Error())
				}
			case 2:
				currentPlugin = 1
				pluginName := strings.TrimSpace(pluginBase02 + "_plugin_v1")
				_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", parameters.DIR_PLUGINS+"/"+pluginName, parameters.DIR_PLUGINS+"/"+pluginBase02+"/"+sourceCode02).CombinedOutput()
				if err != nil {
					fmt.Println("Shared:: Something wrong in generating plugin '"+pluginName+"'")
					os.Stderr.WriteString(err.Error())
				}
			}
			time.Sleep(parameters.INJECTION_TIME * time.Second)
		}
	}
}


func confToGoType(tConf string) string {
	foundType := false
	tGo := ""
	for t := range library.BehaviourLibrary {
		if strings.Contains(t, tConf) {
			tGo = t
			foundType = true
		}
	}

	if !foundType {
		fmt.Println("Type '" + tConf + "' NOT FOUND in Library")
		os.Exit(0)
	}
	return tGo
}

