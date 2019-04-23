package main

import (
	"newelements/artefacts"
	"fmt"
)

func main() {

	madlFile := artefacts.MADLFile{}
	madlFile.ReadMADL("SenderReceiver.madl")

	madl := artefacts.MADL{}
	madl.GenerateMADL(madlFile)

	fmt.Println(madl.Components)
	fmt.Println(madl.Connectors)
	fmt.Println(madl.Attachments)
}
