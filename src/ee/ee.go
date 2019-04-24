package ee

import (
	"newelements/artefacts"
	"fmt"
	"errors"
)

type EE struct{}

func (ee EE) Deploy(madlFileName string) error {
	r1 := *new(error)

	madlFile := artefacts.MADLFile{}
	madlFile.ReadMADL(madlFileName)

	// Create MADL Mid
	madl := artefacts.MADL{}
	err := madl.Create(madlFile)
	if err != nil {
		r1 = errors.New("EE:: "+err.Error())
	}

	// Create MADL EE
	madlEE,err := madl.CreateEE()
	if err != nil {
		r1 = errors.New("EE:: "+err.Error())
	}

	// Generate MADL Go Mid
	madlGo := artefacts.MADLGO{}
	err = madlGo.Create(madl)
	if err != nil {
		r1 = errors.New("EE:: "+err.Error())
	}

	// Generate MADL Go EE
	madlEEGo := artefacts.MADLGO{}
	err = madlEEGo.Create(madlEE)
	if err != nil {
		r1 = errors.New("EE:: "+err.Error())
	}

	// Generate CSP Mid

	// Generate CSP EE

	fmt.Println(madlGo)
	fmt.Println(madlEE.ConfigurationName)

	return r1
}
