package openhab

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/creativeprojects/gopenhab/event"
)

func preventPanic() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "*****************\n")
		fmt.Fprintf(os.Stderr, "***** PANIC *****\n")
		fmt.Fprintf(os.Stderr, "*****************\n\n")
		fmt.Fprintf(os.Stderr, "%s\n\n", r)
		fmt.Fprintf(os.Stderr, "Stack trace:\n\n%s\n", getStack(3)) // skip calls to getStack - preventPanic - panic
		fmt.Fprintf(os.Stderr, "*****************\n\n")
	}
}

func preventRulePanic(client *Client, ruleData RuleData, e event.Event) {
	if r := recover(); r != nil {
		now := time.Now()
		message := fmt.Sprintf("%v", r)
		fmt.Fprintf(os.Stderr, "*****************\n")
		fmt.Fprintf(os.Stderr, "***** PANIC *****\n")
		fmt.Fprintf(os.Stderr, "*****************\n\n")
		fmt.Fprintf(os.Stderr, "Timestamp:   %s\n", now.Format(time.RFC1123))
		fmt.Fprintf(os.Stderr, "Rule ID:     %s\n", ruleData.ID)
		fmt.Fprintf(os.Stderr, "Rule Name:   %s\n", ruleData.Name)
		fmt.Fprintf(os.Stderr, "Description: %s\n", ruleData.Description)
		fmt.Fprintf(os.Stderr, "Event:       %s\n", e.String())
		fmt.Fprintf(os.Stderr, "Message:     %s\n\n", message)
		fmt.Fprintf(os.Stderr, "Stack trace:\n*****************\n%s\n", getStack(3)) // skip calls to getStack - preventRulePanic - panic
		fmt.Fprintf(os.Stderr, "*****************\n\n")
		client.userEventBus.Publish(event.NewRulePanicEvent(message, ruleData.ID, ruleData.Name, ruleData.Description, e.String(), now))
	}
}

// getStack returns a simplified stack trace
func getStack(skip int) string {
	stack := ""
	pc := make([]uintptr, 20)        // stack of 20 traces max
	n := runtime.Callers(skip+1, pc) // skip call to runtime.Callers
	if n == 0 {
		return ""
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	// Loop to get frames.
	// A fixed number of PCs can expand to an indefinite number of Frames.
	for {
		frame, more := frames.Next()

		if frame.Function == "runtime.main" {
			// last 2 traces are inside the go bootstrap
			break
		}

		stack += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)

		if !more {
			break
		}
	}
	if stack == "" {
		// the stack trace is suspiciously empty, let's try another way instead
		return string(debug.Stack())
	}
	return stack
}
