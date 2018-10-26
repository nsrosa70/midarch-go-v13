package parameters

import "time"

// Dirs
//const BASE_DIR  = "/go/midarch-go"  // docker
const BASE_DIR = "/Users/nsr/Dropbox/go/midarch-go-v11"
const DIR_PLUGINS = BASE_DIR + "/src/plugins"
const DIR_CSP = BASE_DIR + "/src/cspspecs"
const DIR_SOURCE = BASE_DIR
const DIR_CONF = BASE_DIR + "/src/apps/conf"
const DIR_GO = "/usr/local/go/bin"
const DIR_FDR = "/Volumes/Macintosh HD/Applications/FDR4-2.app/Contents/MacOS"
const DIR_CSPARSER = "/Users/nsr/Dropbox/go/software-architecture-v10/src/verificationtools/cspdot/csparser"

// Ports
const NAMING_PORT = 4040
const CALCULATOR_PORT = 2020
const FIBONACCI_PORT = 2030
const QUEUEING_PORT = 2040

//
const JAVA_COMMAND = "java"
const JAR_COMMAND = "-jar"

var SetOfPorts = map[string]int{
	"NAMING_PORT":     NAMING_PORT,
	"CALCULATOR_PORT": CALCULATOR_PORT,
	"FIBONACCI_PORT":  FIBONACCI_PORT,
	"QUEUEING_PORT":   QUEUEING_PORT}

const NO_CHANGE = 0
const REACTIVE_CHANGE = 1
const EVOLUTIVE_CHANGE = 2
const PROACTIVE_CHANGE = 3

const CHAN_BUFFER_SIZE = 1
const QUEUE_SIZE = 100

//const PLUGIN_BASE_NAME  = "calculatorinvoker"
const PLUGIN_BASE_NAME = "fibonacciinvoker"
const GRAPH_SIZE = 30

var IS_EVOLUTIVE  = false
var IS_CORRECTIVE = false
var IS_PROACTIVE  = false

var MONITOR_TIME time.Duration   // seconds
var INJECTION_TIME time.Duration // seconds
var REQUEST_TIME time.Duration   // milliseconds
var STRATEGY = 0             // 1 - no change 2 - change once 3 - change same plugin 4 - alternate plugins
var SAMPLE_SIZE = 0
var NAMING_HOST = ""
var QUEUEING_HOST = ""
