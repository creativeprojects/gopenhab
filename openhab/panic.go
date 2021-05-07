package openhab

import (
	"fmt"
	"os"
	"runtime/debug"
)

func preventPanic() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "*****************\n")
		fmt.Fprintf(os.Stderr, "***** PANIC *****\n")
		fmt.Fprintf(os.Stderr, "*****************\n\n")
		fmt.Fprintf(os.Stderr, "%s\n\n", r)
		fmt.Fprintf(os.Stderr, "Stack trace:\n\n%s\n", debug.Stack())
		fmt.Fprintf(os.Stderr, "*****************\n\n")
	}
}
