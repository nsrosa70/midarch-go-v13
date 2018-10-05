package library

import (
	"apps/calculator/calculatorclientproxy"
	"framework/components/requestor"
	"framework/connectors"
	"framework/components/srh"
	"framework/components/naming/namingclientproxy"
	"framework/components/crh"
	"apps/calculator/calculatorinvoker"
	"apps/fibonacci/fibonacciclientproxy"
	"apps/fibonacci/fibonacciinvoker"
	"framework/components/queueing/queueingclientproxy"
	"framework/components/naming/naminginvoker"
	"framework/components/queueing/queueinginvoker"
	"framework/components/sender"
	"framework/components/receiver"
	"framework/components/clientX"
	"framework/components/server"
)

type Record struct {
	RBD   string
	PRISM string
	CSP   string
	Go    interface{}
	PetriNet string
}

var Repository = map[string]Record{
	"calculatorclientproxy.CalculatorClientProxy": Record{RBD: "TODO", PRISM: "TODO", Go: calculatorclientproxy.CalculatorClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"calculatorinvoker.CalculatorInvoker": Record{RBD: "TODO", PRISM: "TODO", Go: calculatorinvoker.CalculatorInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"fibonacciclientproxy.FibonacciClientProxy": Record{RBD: "TODO", PRISM: "TODO", Go: fibonacciclientproxy.FibonacciClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"fibonacciinvoker.FibonacciInvoker": Record{RBD: "TODO", PRISM: "TODO", Go: fibonacciinvoker.FibonacciInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"requestor.Requestor": Record{RBD: "TODO", PRISM: "TODO", Go: requestor.Requestor{}, CSP: "B = InvP.e1 -> I_PosInvP -> InvR.e2 -> TerR.e2 -> I_PosTerR -> TerP.e1 -> B"},
	"connectors.RequestReply": Record{RBD: "TODO", PRISM: "TODO", Go: connectors.RequestReply{}, CSP: "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B"},
	"connectors.NTo1": Record{RBD: "TODO", PRISM: "TODO", Go: connectors.NTo1{}, CSP: "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B [] InvP.e3 -> InvR.e2 -> TerR.e2 -> TerP.e3 -> B"},
	"connectors.OneWay": Record{RBD: "TODO", PRISM: "TODO", Go: connectors.OneWay{}, CSP: "B = InvP.e1 -> InvR.e2 -> B"},
	"sender.Sender{}":Record{RBD: "TODO", PRISM: "TODO", Go: sender.Sender{}, CSP: "B = InvR.e1 -> B"},
	"receiver.Receiver{}":Record{RBD: "TODO", PRISM: "TODO", Go: receiver.Receiver{}, CSP: "B = InvP.e1 -> I_PosInvP -> B"},
	"clientX.ClientX{}":Record{RBD: "TODO", PRISM: "TODO", Go: clientX.ClientX{}, CSP: "B = InvR.e1 -> TerR.e1 -> B"},
	"server.Server{}":Record{RBD: "TODO", PRISM: "TODO", Go: server.Server{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"naminginvoker.NamingInvoker": Record{RBD: "TODO", PRISM: "TODO", Go: naminginvoker.NamingInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"queueinginvoker.QueueingInvoker": Record{RBD: "TODO", PRISM: "TODO", Go: queueinginvoker.QueueingInvoker{}, CSP: "B = InvP.e1 -> I_PosInvP -> TerP.e1 -> B"},
	"namingclientproxy.NamingClientProxy": Record{RBD: "TODO", PRISM: "TODO", Go: namingclientproxy.NamingClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"queueingclientproxy.QueueingClientProxy": Record{RBD: "TODO", PRISM: "TODO", Go: queueingclientproxy.QueueingClientProxy{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"srh.SRH": Record{RBD: "TODO", PRISM: "TODO", Go: srh.SRH{}, CSP: "B = I_PreInvR -> InvR.e1 -> TerR.e1 -> I_PosTerR -> B"},
	"crh.CRH": Record{RBD: "TODO", PRISM: "TODO", Go: crh.CRH{}, CSP: "B = InvP.e1 -> I_PosInvP -> I_PreTerP -> TerP.e1 -> B"}}