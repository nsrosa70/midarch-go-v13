package shared

import (
	"reflect"
	"framework/messages"
	"strings"
	"shared/parameters"
	"strconv"
	"time"
	"fmt"
	"os"
)

const PREFIX_ACTION = "->"

//const CHOICE = "[]"
const PREFIX_INTERNAL_ACTION = "I_"
const INVP = "InvP"
const TERP = "TerP"
const INVR = "InvR"
const TERR = "TerR"
const EVOLUTIVE = "EVOLUTIVE"
const CORRECTIVE = "REACTIVE"
const PROACTIVE = "PROACTIVE"
const NONE = "NONE"

type Invocation struct {
	Method  reflect.Value
	InArgs  []reflect.Value
	OutArgs [] reflect.Value
}

// MAPE-K Types
type MonitoredCorrectiveData string   // used in channel Monitor -> Analyser (Corrective)
type MonitoredEvolutiveData []string  // used in channel Monitor -> Analyser (Evolutive)
type MonitoredProactiveData [] string // used in channel Monitor -> Analyser (Proactive)

type AnalysisResult struct {
	Result   interface{}
	Analysis int
}

type AdaptationPlan struct{
	Plan string
}

var ValidActions = map[string]bool{
	INVP: true,
	TERP: true,
	INVR: true,
	TERR: true}

func IsInternal(action string) bool {
	r := false
	if len(action) >= 2 {
		if action[0:2] == PREFIX_INTERNAL_ACTION {
			r = true
		}
	} else {
		r = false
	}
	return r
}

func IsExternal(action string) bool {
	r := false

	if len(action) >= 2 {
		action := strings.TrimSpace(action)
		if strings.Contains(action,"."){
			action = action[:strings.Index(action,".")]
		}

		if action == INVP || action == TERP || action == INVR || action == TERR {
			r = true
		} else {
			r = false
		}
	} else {
		r = false
	}
	return r
}

func IsAction(action string) bool {
	r := false

	action = strings.TrimSpace(action)
	if len(action) > 2 {
		if strings.Contains(action, INVP) || strings.Contains(action, TERP) || strings.Contains(action, INVR) || strings.Contains(action, TERR) {
			r = true
		} else {
			r = false
		}
	}
	return r
}

type ParMapActions struct {
	F1 func(*chan messages.SAMessage, *messages.SAMessage)      // External action
	F2 func(any interface{}, name string, args ... interface{}) // Internal action
	P1 *messages.SAMessage
	P2 *chan messages.SAMessage
	P3 interface{}
	P4 string
}

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func IsToElement(action string) bool {
	if action[:2] == PREFIX_INTERNAL_ACTION || action[:4] == INVP || action[:4] == TERR {
		return true
	} else { // TerP and InvR
		return false
	}
}

func LoadParameters(args []string) {
	for i := range args {
		variable := strings.Split(args[i], "=")
		switch strings.TrimSpace(variable[0]) {
		case "SAMPLE_SIZE":
			parameters.SAMPLE_SIZE, _ = strconv.Atoi(variable[1])
		case "REQUEST_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			parameters.REQUEST_TIME = temp2
		case "INJECTION_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			parameters.INJECTION_TIME = temp2
		case "MONITOR_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			parameters.MONITOR_TIME = temp2
		case "STRATEGY":
			parameters.STRATEGY, _ = strconv.Atoi(variable[1])
		case "NAMING_HOST":
			parameters.NAMING_HOST = variable[1]
		case "QUEUEING_HOST":
			parameters.QUEUEING_HOST = variable[1]
		default:
			fmt.Println("Shared:: Parameter '" + variable[0] + "' does not exist")
			os.Exit(0)
		}
	}
}

func ShowExecutionParameters(s bool) {
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
		fmt.Println("Adaptability  ")
		fmt.Println("Corrective        : " + strconv.FormatBool(parameters.IS_CORRECTIVE))
		fmt.Println("Evolutive         : " + strconv.FormatBool(parameters.IS_EVOLUTIVE))
		fmt.Println("Proactive         : " + strconv.FormatBool(parameters.IS_PROACTIVE))
		fmt.Println("Monitor Time (s)  : " + (parameters.MONITOR_TIME * time.Second).String())
		fmt.Println("Injection Time (s): " + (parameters.INJECTION_TIME * time.Second).String())
		fmt.Println("Request Time (ms) : " + parameters.REQUEST_TIME.String())
		fmt.Println("Strategy (0-NOT DEFINED 1-No change 2-Change once 3-change same plugin 4-alternate plugins): " + strconv.Itoa(parameters.STRATEGY))
		fmt.Println("******************************************")
	}
}

func Log(args ...string) {
	if strings.Contains(args[1], "Proxy") || strings.Contains(args[1], "XXX") {
		fmt.Println(args[0] + ":" + args[1] + ":" + args[2] + ":" + args[3])
	}
}

func IsConnectorInAttachments(atts []string, t string) bool {
	foundConnector := false

	for a := range atts {
		att := strings.Split(atts[a], ",")
		if (strings.TrimSpace(att[1]) == t) {
			foundConnector = true
		}
	}

	return foundConnector
}

func IsComponentInAttachments(atts []string, c string) bool {
	foundComponent := false

	for a := range atts {
		att := strings.Split(atts[a], ",")
		if (strings.TrimSpace(att[0]) == c || strings.TrimSpace(att[2]) == c) {
			foundComponent = true
		}
	}

	return foundComponent
}

func IsInConnectors(conns map[string]string, t string) bool {
	foundConnector := false

	for i := range conns {
		if t == i {
			foundConnector = true
			break
		}
	}
	return foundConnector
}

/*func RenameInternalChannels(b string, id string) string {
	tokens := strings.Split(b, " ")

	for t := range tokens {
		if IsInternal(tokens[t]) {
			b = strings.Replace(b, tokens[t], tokens[t]+"_"+id, 99)
		}
	}
	return b
}
*/
type QueueingInvocation struct {
	Op   string
	Args []interface{}
}

type QueueingTermination struct {
	R interface{}
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
