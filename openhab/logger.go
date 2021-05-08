package openhab

import "log"

// Logger interface used by gopenhab so you can bring your own logger
type Logger interface {
	Printf(format string, v ...interface{})
}

type emptyLogger struct{}

func (l emptyLogger) Printf(format string, v ...interface{}) {}

var (
	debuglog Logger = &emptyLogger{}
	errorlog Logger = log.Default()
)

func SetErrorLog(logger Logger) {
	errorlog = logger
}

func SetDebugLog(logger Logger) {
	debuglog = logger
}
