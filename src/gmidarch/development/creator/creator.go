package creator

import (
	"errors"
	"gmidarch/development/artefacts/madl"
)

type Creator struct{}

func (Creator) Create(madlFileName string) (madl.MADLGo, madl.MADLGo, error) {
	r1 := madl.MADLGo{}
	r2 := madl.MADLGo{}
	r3 := *new(error)

	madlFile := madl.MADLFile{}
	madlFile.Read(madlFileName)

	// Create MADL Mid
	madlMid := madl.MADL{}
	r3 = madlMid.Create(madlFile)
	if r3 != nil {
		r3 = errors.New("Creator:: " + r3.Error())
		return r1, r2, r3
	}

	// Create MADL EE
	madlEE, r3 := madlMid.CreateEE(madlMid.Adaptability)
	if r3 != nil {
		r3 = errors.New("Creator:: " + r3.Error())
		return r1,r2,r3
	}

	// Create MADL Go Mid
	madlMidGo := madl.MADLGo{}
	r3 = madlMidGo.Create(madlMid)
	if r3 != nil {
		r3 = errors.New("Creator:: " + r3.Error())
		return r1,r2,r3
	}
	r1 = madlMidGo

	// Create MADL Go EE
	madlEEGo := madl.MADLGo{}
	r3 = madlEEGo.Create(madlEE)
	if r3 != nil {
		r3 = errors.New("Creator:: " + r3.Error())
		return r1,r2,r3
	}
	r2 = madlEEGo

	//madlEE.Print()
	//os.Exit(0)

	return r1, r2, r3
}
