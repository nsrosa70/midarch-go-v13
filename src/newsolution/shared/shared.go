
package shared

const EXECUTE_FOREVER = true

type Request struct{
	Op string
	Args []interface{}
}

type Invocation struct {
	Host string
	Port int
	Req Request
}