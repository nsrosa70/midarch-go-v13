package shared

import (
	"reflect"
	"framework/message"
	"strings"
	"os"
	"os/exec"
	"time"
	"shared/parameters"
	"framework/library"
	"fmt"
	"strconv"
)

const PREFIX_ACTION = "->"
const CHOICE = "[]"

type Invocation struct {
	Method  reflect.Value
	InArgs  []reflect.Value
	OutArgs [] reflect.Value
}

type MonitoredCorrectiveData string  // used in channel Monitor -> Analyser (Corrective)
type MonitoredEvolutiveData []string // used in channel Monitor -> Analyser (Evolutive)

type AnalysisResult struct {
	// used in channel Analyser -> Planner
	Result   interface{}
	Analysis int
}

var ValidActions = map[string]bool{
	"InvP": true,
	"TerP": true,
	"InvR": true,
	"TerR": true}

func IsInternal(action string) bool {
	return action[0:2] == "I_"
}

func ToActions(behaviour string) [] string {
	// B = InvP -> B [] InvP -> B
	var actions []string

	if !strings.Contains(behaviour, "[]") {
		behaviourTemp := strings.Split(behaviour, "=")
		behaviour = behaviourTemp[1][0:strings.LastIndex(behaviourTemp[1], "->")]
		actions = strings.Split(behaviour, PREFIX_ACTION)
		if len(actions) == 0 {
			actions[0] = strings.TrimSpace(behaviour)
		} else {
			for i := range actions {
				actions[i] = strings.TrimSpace(actions[i])
			}
		}
	} else {
		behaviourTemp := strings.Split(behaviour, "=")
		branches := strings.Split(behaviourTemp[1], "[]")
		idx := 0
		for i := range branches {
			actionsTemp := strings.Split(branches[i], PREFIX_ACTION)
			for j := range actionsTemp {
				action := strings.TrimSpace(actionsTemp[j])
				if action != "B" && action != "" {
					actions = append(actions, strings.TrimSpace(action))
					idx++
				}
			}
		}
	}
	return actions
}

type ParMapActions struct {
	F1 func(*chan message.Message, *message.Message)            // External action
	F2 func(any interface{}, name string, args ... interface{}) // Internal action
	P1 *message.Message
	P2 *chan message.Message
	P3 interface{}
	P4 string
}

type SubMessage struct {
	I int
}
type MyMessage struct {
	Name string
	Age  int
	X    interface{}
}

func Invoke(any interface{}, name string, args ... interface{}) {
	inputs := make([]reflect.Value, len(args))

	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)

	inputs = nil
	return
}

func SelectChan(action string, id string, channs map[string]chan message.Message, elemMaps map[string]string) chan message.Message {

	p1 := action[0:strings.Index(action, ".")]
	p2 := action[strings.Index(action, ".")+1 : len(action)]

	keyMaps := id + "." + p2
	keyChannel := id + "." + p1 + "." + elemMaps[keyMaps]

	return channs[keyChannel]
}

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

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func ProcessOSArguments(args []string){
	for i:= range args{
		variable := strings.Split(args[i],"=")
		switch strings.TrimSpace(variable[0]) {
		case "SAMPLE_SIZE":
			parameters.SAMPLE_SIZE,_ = strconv.Atoi(variable[1])
		case "REQUEST_TIME":
			temp1,_ := strconv.Atoi(variable[1])
			temp2   := time.Duration(temp1)
			parameters.REQUEST_TIME = temp2
		case "INJECTION_TIME":
			temp1,_ := strconv.Atoi(variable[1])
			temp2   := time.Duration(temp1)
			parameters.INJECTION_TIME = temp2
		case "MONITOR_TIME":
			temp1,_ := strconv.Atoi(variable[1])
			temp2   := time.Duration(temp1)
			parameters.MONITOR_TIME = temp2
		case "STRATEGY":
			parameters.STRATEGY,_ = strconv.Atoi(variable[1])
		case "IS_ADAPTIVE":
			parameters.IS_ADAPTIVE,_ = strconv.ParseBool(variable[1])
		case "NAMING_HOST":
			parameters.NAMING_HOST = variable[1]
		default:
			fmt.Println("Argument '"+variable[0]+"' does not exist")
			os.Exit(0)
		}
	}
}