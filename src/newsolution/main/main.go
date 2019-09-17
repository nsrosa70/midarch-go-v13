package main

import (
	"newsolution/components"
	"newsolution/engine"
	"fmt"
)

func main(){
	c := components.NewClient()
	e := engine.Engine{}
	u := components.NewUnit(c)

	go e.Execute(u)

	fmt.Scanln()
}
