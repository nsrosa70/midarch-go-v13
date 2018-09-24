package library

import (
	"apps/calculator/calculatorclientproxy"
	"framework/components/requestor"
	"framework/connectors"
	"framework/components/invoker/naminginvoker"
	"framework/components/srh"
	"framework/components/naming/namingclientproxy"
	"framework/components/naming/naming"
	"framework/components/crh"
	"apps/calculator/calculatorinvoker"
	"apps/fibonacci/fibonacciclientproxy"
	"apps/fibonacci/fibonacciinvoker"
	"framework/components/queueing"
)

type Record struct {
	RBD   string
	PRISM string
	CSP   string
	Go    interface{}
}

var Repository = map[string]Record{
	"calculatorclientproxy.CalculatorClientProxy": Record{RBD: "TODO", PRISM: "TODO", Go: calculatorclientproxy.CalculatorClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"calculatorinvoker.CalculatorInvoker": Record{RBD: "TODO", PRISM: "TODO", Go: calculatorinvoker.CalculatorInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"fibonacciclientproxy.FibonacciClientProxy": Record{RBD: "TODO", PRISM: "TODO", Go: fibonacciclientproxy.FibonacciClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"fibonacciinvoker.FibonacciInvoker": Record{RBD: "TODO", PRISM: "TODO", Go: fibonacciinvoker.FibonacciInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"requestor.Requestor": Record{RBD: "TODO", PRISM: "TODO", Go: requestor.Requestor{}, CSP: "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> TerR.e2 -> I_PosTerR -> TerP.e1 -> B"},
	"connectors.RequestReply": Record{RBD: "TODO", PRISM: "TODO", Go: connectors.RequestReply{}, CSP: "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B"},
	"connectors.NTo1": Record{RBD: "TODO", PRISM: "TODO", Go: connectors.NTo1{}, CSP: "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B [] InvP.e3 -> InvR.e2 -> TerR.e2 -> TerP.e3 -> B"},
	"naming.NamingService": Record{RBD: "TODO", PRISM: "TODO", Go: naming.NamingService{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"naminginvoker.NamingInvoker": Record{RBD: "TODO", PRISM: "TODO", Go: naminginvoker.NamingInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> TerR.e2 -> I_PosTerR -> TerP.e1 -> B"},
	"namingclientproxy.NamingClientProxy": Record{RBD: "TODO", PRISM: "TODO", Go: namingclientproxy.NamingClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"srh.SRH": Record{RBD: "TODO", PRISM: "TODO", Go: srh.SRH{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"crh.CRH": Record{RBD: "TODO", PRISM: "TODO", Go: crh.CRH{}, CSP: "B = InvP.e1 -> I_PosInvP -> I_PreTerP -> TerP.e1 -> B"},
	"queueing.QueueingService": Record{RBD: "TODO", PRISM: "TODO", Go: queueing.QueueingService{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
    }

//var TypeLibrary = map[string]interface{}{
	//"calculatorclientproxy.CalculatorClientProxy": calculatorclientproxy.CalculatorClientProxy{},
	//"calculatorinvoker.CalculatorInvoker":         calculatorinvoker.CalculatorInvoker{},
//	"fibonacciclientproxy.FibonacciClientProxy":   fibonacciclientproxy.FibonacciClientProxy{},
//	"fibonacciinvoker.FibonacciInvoker":           fibonacciinvoker.FibonacciInvoker{},
//	"requestor.Requestor":                         requestor.Requestor{},
//	"connectors.RequestReply":                     connectors.RequestReply{},
//	"connectors.NTo1":                             connectors.NTo1{},
//	"naming.NamingService":                        naming.NamingService{},
//	"naminginvoker.NamingInvoker":                 naminginvoker.NamingInvoker{},
//	"namingclientproxy.NamingClientProxy":         namingclientproxy.NamingClientProxy{},
//	"srh.SRH":                                     srh.SRH{},
//	"crh.CRH":                                     crh.CRH{},
//	"queueing.QueueingService":                    queueing.QueueingService{}}

//var BehaviourLibrary = map[string]string{
//	"sender.Sender":                               "B = I_PreInvR -> InvR.e1 -> B",
//	"receiver.Receiver":                           "B = InvP.e1 -> I_PosInvP -> B",
//	"connectors.OneWay":                           "B = InvP.e1 -> InvR.e2 -> B",
//	"client.Client":                               "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B",
//	"client.Client1":                              "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B",
//	"client.Client2":                              "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B",
//	"client.Client3":                              "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B",
//	"server.Server":                               "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B",
//	"connectors.RequestReply":                     "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B",
//	"requestor.Requestor":                         "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> TerR.e2 -> I_PosTerR -> TerP.e1 -> B",
//	"crh.CRH":                                     "B = InvP.e1 -> I_PosInvP -> I_PreTerP -> TerP.e1 -> B",
//	"srh.SRH":                                     "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B",
//	"clientproxy.ClientProxy":                     "B = InvR.e1 -> TerR.e1 -> B",
//	"sender.SenderGeneric":                        "B = I_PreInvR -> InvR.e1 -> B",
//	"naming.NamingService":                        "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B",
//	"naminginvoker.NamingInvoker":                 "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> TerR.e2 -> I_PosTerR -> TerP.e1 -> B",
//	"calculatorinvoker.CalculatorInvoker":         "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B",
//	"fibonacciinvoker.FibonacciInvoker":           "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B",
//	"namingclientproxy.NamingClientProxy":         "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B",
//	"calculatorclientproxy.CalculatorClientProxy": "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B",
//	"fibonacciclientproxy.FibonacciClientProxy":   "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B",
//	"invoker.Invoker":                             "B = InvP.e1 -> TerP.e1 -> B",
//	"connectors.NTo1":                             "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B [] InvP.e3 -> InvR.e2 -> TerR.e2 -> TerP.e3 -> B",
//	"client.Client1TCP":                           "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"}
