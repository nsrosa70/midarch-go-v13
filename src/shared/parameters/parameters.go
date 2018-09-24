package parameters

import "time"

// Dirs
//const BASE_DIR    = "/go/midarch-go"
const BASE_DIR    = "/Users/nsr/Dropbox/go/software-architecture-v8"
const DIR_PLUGINS = BASE_DIR + "/src/plugins"
const DIR_CSP     = BASE_DIR+ "/src/cspspecs"
const DIR_SOURCE  = BASE_DIR
const DIR_CONF    = BASE_DIR+"/src/apps/conf"
const DIR_GO      = "/usr/local/go/bin"

// Ports

const NAMING_PORT     = 4040
const CALCULATOR_PORT = 2020
const FIBONACCI_PORT  = 2121

var SetOfPorts = map[string]int{
	"NAMING_PORT"     : NAMING_PORT,
	"CALCULATOR_PORT" : CALCULATOR_PORT,
	"FIBONACCI_PORT"  : FIBONACCI_PORT}

const NO_CHANGE        = 0
const REACTIVE_CHANGE  = 1
const EVOLUTIVE_CHANGE = 2
const PROACTIVE_CHANGE = 3

const CHAN_BUFFER_SIZE = 1

//const PLUGIN_BASE_NAME  = "calculatorinvoker"
const PLUGIN_BASE_NAME    = "fibonacciinvoker"
const GRAPH_SIZE          = 30

var IS_ADAPTIVE = true
var INJECTION_ENABLED = false
var MONITOR_TIME  time.Duration   // seconds
var INJECTION_TIME time.Duration  // seconds
var REQUEST_TIME time.Duration    // milliseconds
var STRATEGY int      = 0   // 1 - no chanvar ge 2 - change once 3 - change same plugin 4 - alternate plugins
var SAMPLE_SIZE int   = 0
var NAMING_HOST    = ""