package factories

import (
	"framework/components/queueing/queueingclientproxy"
	"shared/parameters"
	"framework/components/naming/namingclientproxy"
)

func FactoryQueueing() queueingclientproxy.QueueingClientProxy {
	p := queueingclientproxy.QueueingClientProxy{Host: parameters.QUEUEING_HOST, Port: parameters.QUEUEING_PORT}
	return p
}

func LocateNaming() namingclientproxy.NamingClientProxy {
	p := namingclientproxy.NamingClientProxy{Host:parameters.NAMING_HOST,Port:parameters.NAMING_PORT}
	return p
}