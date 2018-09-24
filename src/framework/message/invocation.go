package message

type Invocation struct {
	Host string
	Port int
	Op string
	Args interface{}
}
