package frontend

import (
	"newsolution/development/artefacts/madl"
	"newsolution/execution/environment/ee"
)

type FrontEnd struct {}

func (f FrontEnd) Deploy(file string) error {
	r := *new(error)

	// Read madl and generate architectural artefacts
	madlTemp := madl.MADL{}
	madlTemp.Read(file)

	// Invoke execution environment
	eeTemp := ee.NewEE(madlTemp)
	eeTemp.Start()

	return r
}
