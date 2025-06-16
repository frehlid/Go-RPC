package helpers

import (
	"fmt"
)

const VERBOSE = false

func CheckForError(err error) {
	if err != nil {
		panic(err)
	}
}

func VerbosePrint(message ...any) {
	if VERBOSE {
		fmt.Println(message)
	}
}
