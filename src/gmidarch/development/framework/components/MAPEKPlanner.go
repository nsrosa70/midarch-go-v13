package components

import (
	"gmidarch/development/framework/messages"
	"gmidarch/development/framework/element"
)

type MAPEKPlanner struct{}

type MAPEKPlannerInfo struct{
	ConfId string
	Components []element.ElementGo
}

func (MAPEKPlanner) I_Plan(msg *messages.SAMessage, info *interface{}, r *bool) {

}