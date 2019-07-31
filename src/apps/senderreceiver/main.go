package main

import (
	"gmidarch/shared/shared"
	"gmidarch/development/frontend"
	"fmt"
)

func main(){

	err := frontend.FrontEnd{}.Deploy("SenderReceiver.madl")
	shared.CheckError(err,"MAIN")

	fmt.Scanln()
}
