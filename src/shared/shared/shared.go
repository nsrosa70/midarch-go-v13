package shared

import (
	"reflect"
	"framework/message"
	"strings"
	"os"
	"time"
	"shared/parameters"
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
	P5 *chan string // REMOVE
	P6 *string      // REMOVE
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

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func IsToElement(action string) bool {
	if action[:2] == "I_" || action[:4] == "InvP" || action[:4] == "TerR" {
		return true
	} else { // TerP and InvR
		return false
	}
}


func LoadParameters(args []string){
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
		case "QUEUEING_HOST":
			parameters.QUEUEING_HOST = variable[1]
		default:
			fmt.Println("Shared:: Parameter '"+variable[0]+"' does not exist")
			os.Exit(0)
		}
	}
}

func ShowExecutionParameters(s bool){
	if s {
		fmt.Println("******************************************")
		fmt.Println("Sample size                : " + strconv.Itoa(parameters.SAMPLE_SIZE))
		fmt.Println("Direrctory of base code    : " + parameters.BASE_DIR)
		fmt.Println("Directory of plugins       : " + parameters.DIR_PLUGINS)
		fmt.Println("Directory of CSP specs     : " + parameters.DIR_CSP)
		fmt.Println("Directory of Configurations: " + parameters.DIR_CONF)
		fmt.Println("Directory of Go compiler   : " + parameters.DIR_GO)
		fmt.Println("Directory of FDR           : " + parameters.DIR_FDR)
		fmt.Println("------------------------------------------")
		fmt.Println("Naming Host     : " + parameters.NAMING_HOST)
		fmt.Println("Naming Port     : " + strconv.Itoa(parameters.NAMING_PORT))
		fmt.Println("Calculator Port : " + strconv.Itoa(parameters.CALCULATOR_PORT))
		fmt.Println("Fibonacci Port  : " + strconv.Itoa(parameters.FIBONACCI_PORT))
		fmt.Println("Queueing Port   : " + strconv.Itoa(parameters.QUEUEING_PORT))
		fmt.Println("------------------------------------------")
		fmt.Println("Plugin Base Name: " + parameters.PLUGIN_BASE_NAME)
		fmt.Println("Max Graph Size  : " + strconv.Itoa(parameters.GRAPH_SIZE))
		fmt.Println("------------------------------------------")
		fmt.Println("Adaptive          : " + strconv.FormatBool(parameters.IS_ADAPTIVE))
		//fmt.Println("Injection enabled : " + strconv.FormatBool(parameters.INJECTION_ENABLED))
		fmt.Println("Monitor Time (s)  : " + (parameters.MONITOR_TIME*time.Second).String())
		fmt.Println("Injection Time (s): " + (parameters.INJECTION_TIME*time.Second).String())
		fmt.Println("Request Time (ms) : " + parameters.REQUEST_TIME.String())
		fmt.Println("Strategy (0-NOT DEFINED 1-No change 2-Change once 3-change same plugin 4-alternate plugins): "+strconv.Itoa(parameters.STRATEGY))
		fmt.Println("******************************************")
	}
}

