package openhabtest

type Logger interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

type dummyLogger struct{}

func (l dummyLogger) Log(args ...interface{})                 {}
func (l dummyLogger) Logf(format string, args ...interface{}) {}
