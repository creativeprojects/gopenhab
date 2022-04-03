package openhab

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func preventPanic() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "*****************\n")
		fmt.Fprintf(os.Stderr, "***** PANIC *****\n")
		fmt.Fprintf(os.Stderr, "*****************\n\n")
		fmt.Fprintf(os.Stderr, "%s\n\n", r)
		fmt.Fprintf(os.Stderr, "Stack trace:\n\n%s\n", getStack())
		fmt.Fprintf(os.Stderr, "*****************\n\n")
	}
}

func getStack() string {
	stack := ""
	pc := make([]uintptr, 10)
	n := runtime.Callers(4, pc)
	if n == 0 {
		return ""
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	// Loop to get frames.
	// A fixed number of PCs can expand to an indefinite number of Frames.
	for {
		frame, more := frames.Next()

		// stop when we get to the runtime bootstrap
		if strings.Contains(frame.File, "runtime/") {
			break
		}
		stack += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)

		if !more {
			break
		}
	}
	return stack
}
