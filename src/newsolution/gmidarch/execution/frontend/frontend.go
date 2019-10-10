package frontend

import (
	"newsolution/gmidarch/execution/loader"
	"newsolution/gmidarch/execution/ee"
	"newsolution/gmidarch/execution/creator"
	"fmt"
)

type FrontEnd struct {}

func (f FrontEnd) Deploy(file string) {

	// Read madl and generate architectural artefacts
	l := loader.Loader{}
	mapp := l.Load(file)

	// Create architecture of the execution environment
	kindOfAdaptability := make([]string,1)
	kindOfAdaptability[0]= "NONE"
	crt := creator.Creator{}
	mee := crt.CreateEE(mapp,kindOfAdaptability)

	fmt.Println(mee)
	// Generate CSPs
	// g := generator.Generator{}
	// cspapp := g.Generate(mapp)
	// cspee := g.Generate(mee)

	// Check CSP
	//checker := checker.Checker{}
	//

	// Start execution environment
	eeTemp := ee.NewEE(mapp)
	eeTemp.Start()

}
