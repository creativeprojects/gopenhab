package openhab

type Logger interface {
	Printf(format string, v ...interface{})
}

type emptyLogger struct{}

func (l emptyLogger) Printf(format string, v ...interface{}) {}

var (
	log Logger = &emptyLogger{}
)

func SetLogger(logger Logger) {
	log = logger
}
