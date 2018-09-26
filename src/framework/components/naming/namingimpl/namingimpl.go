package namingimpl

import "framework/components/naming/ior"

type NamingImpl struct{}

var Repo = map[string]ior.IOR{}

func (NamingImpl) Lookup(s string) ior.IOR {
	return Repo[s]
}

func (NamingImpl) List() []string{
	keys := make([]string, 0, len(Repo))
	for k := range Repo {
		keys = append(keys, k)
	}
	return keys
}

func (n NamingImpl) Register(serviceName string, ior ior.IOR) bool{
	if _, ok := Repo[serviceName]; ok {
		return false
	} else {
		Repo[serviceName] = ior
		return true
	}
}

