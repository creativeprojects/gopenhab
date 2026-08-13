package openhabtest

type Logger interface {
	Log(args ...any)
	Logf(format string, args ...any)
}

type dummyLogger struct{}

func (l dummyLogger) Log(args ...any)                 {}
func (l dummyLogger) Logf(format string, args ...any) {}
