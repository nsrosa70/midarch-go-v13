package architectural

import (
	"newsolution/development/components"
	"newsolution/development/connectors"
)

type Record struct {
	Type interface{}
	CSP string
	X interface{}
}

type ArchitecturalRepository struct {
	Library map[string]Record
}

func (al *ArchitecturalRepository) Load() error {
	r1 := *new(error)

	al.Library = make(map[string]Record)

	// load
	al.Library["Oneway"] = Record{Type:connectors.NewOneway(),X:connectors.Oneway{},CSP:"TODO"}
	al.Library["Requestreply"] = Record{Type:connectors.NewRequestReply(),CSP:""}
	al.Library["Sender"] = Record{Type:components.NewSender(),CSP:""}
	al.Library["Receiver"] = Record{Type:components.NewReceiver(),CSP:""}
	al.Library["Client"] = Record{Type:components.NewClient(),CSP:""}
	al.Library["Server"] = Record{Type:components.NewServer(),CSP:""}
	al.Library["Calculatorproxy"] = Record{Type:components.NewCalculatorProxy(),CSP:""}
	al.Library["Marshaller"] = Record{Type:components.NewMarshaller(),CSP:""}
	al.Library["Requestor"] = Record{Type:components.NewRequestor(),CSP:""}
	al.Library["CRH"] = Record{Type:components.NewCRH(),CSP:""}
	al.Library["SRH"] = Record{Type:components.NewSRH(),CSP:""}
	al.Library["Calculatorserver"] = Record{Type:components.Newcalculatorserver(),CSP:""}
	al.Library["Calculatorinvoker"] = Record{Type:components.NewCalculatorinvoker(),CSP:""}
	al.Library["Calculatorclient"] = Record{Type:components.NewCalculatorclient(),CSP:""}

	return r1
}