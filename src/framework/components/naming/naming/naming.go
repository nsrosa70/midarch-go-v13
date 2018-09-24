package naming

import (
	"framework/message"
	"transport/middleware/ior"
	"fmt"
	"framework/components/naming/namingclientproxy"
)

type NamingService struct{}

var Repo = map[string]ior.IOR{}

func (n NamingService) I_PosInvP(msg *message.Message) {

	// recover parameters
	op := msg.Payload.(message.Invocation).Op
	args := msg.Payload.(message.Invocation).Args

	switch op {
	case "register":
		fmt.Println("Naming:: Register")
	    argsX := args.([]interface{})
		fmt.Println(argsX)
	    p1 := argsX[0].(string)
		p2 := argsX[1].(map[string]interface{})
		ior := ior.IOR{Host:p2["Host"].(string),Port:int(p2["Port"].(float64)),Id:int(p2["Id"].(float64)),Proxy:p2["Proxy"].(string)}
		r := n.Register(p1,ior)
		msgRep := message.Message{r}
		*msg = msgRep
	case "lookup":
		fmt.Println("Naming:: Lookup")
		argsX := args.([]interface{})
		p1 := argsX[0].(string)
		r := n.Lookup(p1)
		msgRep := message.Message{r}
		*msg = msgRep
	case "list":
		fmt.Println("Naming:: List")
		r := n.List()
		msgRep := message.Message{r}
		*msg = msgRep
	}
}

// ************* Functional interface ****************//
func (NamingService) Lookup(s string) ior.IOR {
	return Repo[s]
}

func (NamingService) List() []string{
	keys := make([]string, 0, len(Repo))
	for k := range Repo {
		keys = append(keys, k)
	}
	return keys
}

func (n NamingService) Register(serviceName string, ior ior.IOR) bool{
	if _, ok := Repo[serviceName]; ok {
		return false
	} else {
		Repo[serviceName] = ior
		return true
	}
 }

func LocateNaming(host string, port int) namingclientproxy.NamingClientProxy {
	p := namingclientproxy.NamingClientProxy{Host:host,Port:port}
	return p
}