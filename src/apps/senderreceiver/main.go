package main

import (
	"gmidarch/shared/shared"
	"gmidarch/development/frontend"
)

func main(){

	err := frontend.FrontEnd{}.Deploy("SenderReceiver.madl")
	shared.CheckError(err,"MAIN")
}
