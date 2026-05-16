package assert

import (
	"log"
	"os"
	"runtime"
	"runtime/debug"
)

func Assert(condition bool, message string) {
	if condition {
		return
	}

	_, file, line, _ := runtime.Caller(1)

	log.Println("Assertion failed: ", file, "\nLine", line, message, "\nStack Trace\n")
	debug.PrintStack()

	os.Exit(-1)
}
