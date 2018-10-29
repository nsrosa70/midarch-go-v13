package components

import (
	"framework/element"
)

type Naming struct{}

var Repo = map[string]element.IOR{}

func (Naming) Lookup(s string) element.IOR {
	return Repo[s]
}

func (Naming) List() []string{
	keys := make([]string, 0, len(Repo))
	for k := range Repo {
		keys = append(keys, k)
	}
	return keys
}

func (Naming) Register(serviceName string, ior ior.IOR) bool{
	if _, ok := Repo[serviceName]; ok {
		return false
	} else {
		Repo[serviceName] = ior
		return true
	}
}

