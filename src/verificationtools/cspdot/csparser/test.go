package main

import (
	"os/exec"
	"shared/parameters"
	"shared/errors"
)

func main() {

	// createDot("conf.csp", "")

	createDot("conf.csp", "requestor")

}

func createDot(cspfile string, process string) {

	parser := parameters.DIR_CSPARSER + "/" + "CSParser.jar"

	inputFile := parameters.DIR_CSPARSER + "/" + cspfile

	out, err := exec.Command(parameters.JAVA_COMMAND, parameters.JAR_COMMAND, parser, inputFile, process).Output()
	if err != nil {
		myError := errors.MyError{Source: "CSParser", Message: "Problem in creating .dot file"}
		myError.ERROR()
	}

	println(out)
}
