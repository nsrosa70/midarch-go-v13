package creator

import (
	"gmidarch/development/artefacts"
	"errors"
)

type Creator struct{}

func (Creator) Create(madlFileName string) (artefacts.MADLGo, artefacts.MADLGo, error) {
	r1 := artefacts.MADLGo{}
	r2 := artefacts.MADLGo{}
	r3 := *new(error)

	madlFile := artefacts.MADLFile{}
	madlFile.Read(madlFileName)

	// Create MADL Mid
	madlMid := artefacts.MADL{}
	r3 = madlMid.Create(madlFile)
	if r3 != nil {
		r3 = errors.New("Creator:: " + r3.Error())
		return r1, r2, r3
	}

	// Create MADL EE
	madlEE, r3 := madlMid.CreateEE()
	if r3 != nil {
		r3 = errors.New("Creator:: " + r3.Error())
		return r1,r2,r3
	}

	// Generate MADL Go Mid
	madlMidGo := artefacts.MADLGo{}
	r3 = madlMidGo.Create(madlMid)
	if r3 != nil {
		r3 = errors.New("Creator:: " + r3.Error())
		return r1,r2,r3
	}
	r1 = madlMidGo

	// Generate MADL Go EE
	madlEEGo := artefacts.MADLGo{}
	r3 = madlEEGo.Create(madlEE)
	if r3 != nil {
		r3 = errors.New("Creator:: " + r3.Error())
		return r1,r2,r3
	}
	r2 = madlEEGo

	return r1, r2, r3
}
