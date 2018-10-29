package naming

import (
	"framework/element"
)

type NamingService struct{}

var Repo = map[string]element.IOR{}

func (NamingService) Lookup(s string) element.IOR {
	return Repo[s]
}

func (NamingService) List() []string{
	keys := make([]string, 0, len(Repo))
	for k := range Repo {
		keys = append(keys, k)
	}
	return keys
}

func (NamingService) Register(serviceName string, ior element.IOR) bool{
	if _, ok := Repo[serviceName]; ok {
		return false
	} else {
		Repo[serviceName] = ior
		return true
	}
}

