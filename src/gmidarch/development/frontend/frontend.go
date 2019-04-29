package frontend

import (
	"errors"
	"gmidarch/development/framework/manager"
)

type FrontEnd struct {}

func (FrontEnd) Deploy(madlFileName string) error {
	r1 := *new(error)

	manager := manager.Manager{}

	r1 = manager.Invoke(madlFileName)
	if r1 != nil {
		r1 = errors.New("FrontEnd:: "+r1.Error())
	}

	return r1
}
