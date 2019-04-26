package ee

import (
	"newelements/artefacts"
	"errors"
)

type EE struct{}

func (ee EE) Deploy(madlFileName string) error {
	r1 := *new(error)

	madlFile := artefacts.MADLFile{}
	madlFile.Read(madlFileName)

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
	cspMid := artefacts.CSP{}
	err = cspMid.Create(madlGo)
	if err != nil {
		r1 = errors.New("EE:: "+err.Error())
	}
	cspMidFile := cspMid.GenerateFile()
	err = cspMidFile.Save()
	if err != nil {
		r1 = errors.New("EE:: "+err.Error())
	}
	err = cspMidFile.Check()
	if err != nil {
		r1 = errors.New("EE:: "+err.Error())
	}

	// Generate CSP EE
	cspEE := artefacts.CSP{}
	err = cspEE.Create(madlEEGo)
	if err != nil {
		r1 = errors.New("EE:: "+err.Error())
	}
	cspEEFile := cspEE.GenerateFile()
	err = cspEEFile.Save()
	if err != nil {
		r1 = errors.New("EE:: "+err.Error())
	}
	err = cspEEFile.Check()
	if err != nil {
		r1 = errors.New("EE:: "+err.Error())
	}

	// Generate State Machines

	//fmt.Println(madlGo)
	//fmt.Println(madlEE.ConfigurationName)

	return r1
}
