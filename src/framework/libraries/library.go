package libraries

import (
	"fmt"
	"os"
	"framework/components"
	"framework/connectors"
	"framework/components/namingclientproxy"
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
	"namingclientproxy.NamingClientProxy":      Record{RBD: "TODO", PRISM: "TODO", Go: namingclientproxy.NamingClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"components.CalculatorClientProxy":         Record{LTS: "TODO", RBD: "TODO", PRISM: "TODO", Go: components.CalculatorClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"components.CalculatorInvoker":             Record{RBD: "TODO", PRISM: "TODO", Go: components.CalculatorInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"components.FibonacciClientProxy":          Record{RBD: "TODO", PRISM: "TODO", Go: components.FibonacciClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"components.FibonacciInvoker":              Record{RBD: "TODO", PRISM: "TODO", Go: components.FibonacciInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"components.Requestor":                     Record{RBD: "TODO", PRISM: "TODO", Go: components.Requestor{}, CSP: "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> TerR.e2 -> I_PosTerR -> TerP.e1 -> B"},
	"components.Sender":                        Record{RBD: "TODO", PRISM: "TODO", Go: components.Sender{}, CSP: "B = I_PreInvR1 -> InvR.e1 -> B [] I_PreInvR2 -> InvR.e1 -> B"},
	"components.Receiver":                      Record{RBD: "TODO", PRISM: "TODO", Go: components.Receiver{}, CSP: "B = InvP.e1 -> I_PosInvP -> B"},
	"components.NamingInvoker":                 Record{RBD: "TODO", PRISM: "TODO", Go: components.NamingInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"components.NotificationEngine":            Record{RBD: "TODO", PRISM: "TODO", Go: components.NotificationEngine{}, CSP: "B = InvP.e1 -> I_PosInvP -> ( I_Subscribe -> InvR.e2 -> TerR.e2 -> TerP.e1 -> I_GetSubs -> InvR.e2 -> TerR.e2 -> I_GetResSubs -> B [] I_Unsubscribe -> InvR.e2 -> TerR.e2 -> TerP.e1 -> I_GetSubs -> InvR.e2 -> TerR.e2 -> I_GetResSubs -> B [] I_GetSubscribers -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B [] I_Publish -> I_Pub -> TerP.e1 -> I_Notify -> InvR.e3 -> TerR.e3 -> B )"},
	"components.NotificationEngineInvoker":     Record{RBD: "TODO", PRISM: "TODO", Go: components.NotificationEngineInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> TerR.e2 -> I_PosTerR -> TerP.e1 -> B "},
	"components.NotificationEngineClientProxy": Record{RBD: "TODO", PRISM: "TODO", Go: components.NotificationEngineClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"components.SubscriptionManager":           Record{RBD: "TODO", PRISM: "TODO", Go: components.SubscriptionManager{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"components.NotificationConsumer":          Record{RBD: "TODO", PRISM: "TODO", Go: components.NotificationConsumer{}, CSP: "B = InvP.e1 -> I_Notify -> TerP.e1 -> B"},
	"components.SRH":                           Record{RBD: "TODO", PRISM: "TODO", Go: components.SRH{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"components.CRH":                           Record{RBD: "TODO", PRISM: "TODO", Go: components.CRH{}, CSP: "B = InvP.e1 -> I_PosInvP -> I_PreTerP -> TerP.e1 -> B"},
	"components.MAPEKMonitor":                  Record{RBD: "TODO", PRISM: "TODO", Go: components.MAPEKMonitor{}, CSP: "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> B"},
	"components.MAPEKAnalyser":                 Record{RBD: "TODO", PRISM: "TODO", Go: components.MAPEKAnalyser{}, CSP: "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> B"},
	"components.MAPEKPlanner":                  Record{RBD: "TODO", PRISM: "TODO", Go: components.MAPEKPlanner{}, CSP: "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> B"},
	"components.MAPEKExecutor":                 Record{RBD: "TODO", PRISM: "TODO", Go: components.MAPEKExecutor{}, CSP: "B = InvP.e1 -> I_PosInvP -> B"},
	"components.CorrectiveMonitor":             Record{RBD: "TODO", PRISM: "TODO", Go: components.CorrectiveMonitor{}, CSP: "B = I_GenerateEvent -> InvR.e1 -> B"},
	"components.EvolutiveMonitor":              Record{RBD: "TODO", PRISM: "TODO", Go: components.EvolutiveMonitor{}, CSP: "B = I_GenerateEvent -> InvR.e1 -> B"},
	"components.ProactiveMonitor":              Record{RBD: "TODO", PRISM: "TODO", Go: components.ProactiveMonitor{}, CSP: "B = I_GenerateEvent -> InvR.e1 -> B"},
	"connectors.RequestReply":                  Record{RBD: "TODO", PRISM: "TODO", Go: connectors.RequestReply{}, CSP: "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B"},
	"connectors.TwoToOne":                      Record{RBD: "TODO", PRISM: "TODO", Go: connectors.TwoToOne{}, CSP: "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B [] InvP.e3 -> InvR.e2 -> TerR.e2 -> TerP.e3 -> B"},
	"connectors.ThreeToOne":                    Record{RBD: "TODO", PRISM: "TODO", Go: connectors.ThreeToOne{}, CSP: "B = InvP.e1 -> InvR.e2 -> B [] InvP.e3 -> InvR.e2 -> B [] InvP.e4 -> InvR.e2 -> B"},
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
