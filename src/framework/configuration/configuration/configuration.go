package configuration

import (
	"framework/element"
	"framework/configuration/attachments"
)

type Configuration struct {
	Id string
	Components  map[string] element.Element
	Connectors  map[string] element.Element
	Attachments [] attachments.Attachment
}

func (conf *Configuration) AddComp(comp element.Element) {

	if conf.Components == nil{
		conf.Components = make(map[string] element.Element)
		conf.Components[comp.Id] = comp
	} else {
		conf.Components[comp.Id] = comp
	}
}

func (conf *Configuration) AddConn(conn element.Element) {

	if conf.Connectors == nil{
		conf.Connectors = make(map[string] element.Element)
		conf.Connectors [conn.Id] = conn
	} else {
		conf.Connectors[conn.Id] = conn
	}
}

func (conf *Configuration) AddAtt(a attachments.Attachment) {
	conf.Attachments = append(conf.Attachments, a)
}
