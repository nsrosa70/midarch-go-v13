package libraries

import (
	"fmt"
	"os"
	"framework/components"
	"framework/connectors"
)

type Record struct {
	RBD      string
	PRISM    string
	CSP      string
	Go       interface{}
	PetriNet string
	LTS      string // *.dot file name
}

var Repository = map[string]Record{
	"components.CalculatorClientProxy":         Record{LTS: "TODO", RBD: "TODO", PRISM: "TODO", Go: components.CalculatorClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"components.CalculatorInvoker":             Record{RBD: "TODO", PRISM: "TODO", Go: components.CalculatorInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"components.FibonacciClientProxy":          Record{RBD: "TODO", PRISM: "TODO", Go: components.FibonacciClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"components.FibonacciInvoker":              Record{RBD: "TODO", PRISM: "TODO", Go: components.FibonacciInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"components.Requestor":                     Record{RBD: "TODO", PRISM: "TODO", Go: components.Requestor{}, CSP: "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> TerR.e2 -> I_PosTerR -> TerP.e1 -> B"},
	"components.Sender":                        Record{RBD: "TODO", PRISM: "TODO", Go: components.Sender{}, CSP: "B = I_PreInvR -> InvR.e1 -> B"},
	"components.Receiver":                      Record{RBD: "TODO", PRISM: "TODO", Go: components.Receiver{}, CSP: "B = InvP.e1 -> I_PosInvP -> B"},
	"components.NamingInvoker":                 Record{RBD: "TODO", PRISM: "TODO", Go: components.NamingInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"components.NotificationEngine{}":          Record{RBD: "TODO", PRISM: "TODO", Go: components.NotificationEngine{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"components.NotificationEngineInvoker":     Record{RBD: "TODO", PRISM: "TODO", Go: components.NotificationEngineInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> TerR.e2 -> I_PosTerR -> TerP.e1 -> B "},
	"components.NamingClientProxy":             Record{RBD: "TODO", PRISM: "TODO", Go: components.NamingClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"components.NotificationEngineClientProxy": Record{RBD: "TODO", PRISM: "TODO", Go: components.NotificationEngineClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"components.SRH":                           Record{RBD: "TODO", PRISM: "TODO", Go: components.SRH{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"components.CRH":                           Record{RBD: "TODO", PRISM: "TODO", Go: components.CRH{}, CSP: "B = InvP.e1 -> I_PosInvP -> I_PreTerP -> TerP.e1 -> B"},
	"connectors.RequestReply":                  Record{RBD: "TODO", PRISM: "TODO", Go: connectors.RequestReply{}, CSP: "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B"},
	"connectors.NTo1":                          Record{RBD: "TODO", PRISM: "TODO", Go: connectors.NTo1{}, CSP: "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B [] InvP.e3 -> InvR.e2 -> TerR.e2 -> TerP.e3 -> B"},
	"connectors.OneWay":                        Record{RBD: "TODO", PRISM: "TODO", Go: connectors.OneWay{}, CSP: "B = InvP.e1 -> InvR.e2 -> B"},}

func CheckLibrary() bool {
	r := true
	for e := range Repository {
		if Repository[e].CSP == "" {
			fmt.Println("Library:: Behaviour of Record '" + e + "' is INVALID!!")
			os.Exit(0)
		}
	}
	return r
}
