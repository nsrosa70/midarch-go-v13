package architectural

import (
	"newsolution/development/connectors"
	"newsolution/development/components"
)

type Record struct {
	Type interface{}
	CSP string
}

type ArchitecturalRepository struct {
	Library map[string]Record
}

func (al *ArchitecturalRepository) Load() error {
	r1 := *new(error)

	al.Library = make(map[string]Record)

	// load
	al.Library["Oneway"] = Record{Type:connectors.NewOneway(),CSP:""}
	al.Library["Requestreply"] = Record{Type:connectors.NewRequestReply(),CSP:""}
	al.Library["Sender"] = Record{Type:components.NewSender(),CSP:""}
	al.Library["Receiver"] = Record{Type:components.NewReceiver(),CSP:""}
	al.Library["Client"] = Record{Type:components.NewClient(),CSP:""}
	al.Library["Server"] = Record{Type:components.NewServer(),CSP:""}
	al.Library["Calculatorproxy"] = Record{Type:components.NewCalculatorProxy(),CSP:""}
	al.Library["Marshaller"] = Record{Type:components.NewMarshaller(),CSP:""}
	al.Library["Requestor"] = Record{Type:components.NewRequestor(),CSP:""}
	al.Library["CRH"] = Record{Type:components.NewCRH(),CSP:""}

	return r1
}