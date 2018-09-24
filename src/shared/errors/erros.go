package errors

import (
	"fmt"
	"os"
)

type MyError struct {
	Source  string
	Message string
}

func (e MyError) ERROR() {
	fmt.Println("[FATAL ERROR][" + e.Source + "] " + e.Message)
	os.Exit(0)
}
