package openhab

import (
	"io"
	"log"
)

type TelemetryLogger struct {
	logger *log.Logger
}

func NewTelemetryLogger(writer io.Writer) *TelemetryLogger {
	logger := log.New(writer, "", log.LstdFlags)
	return &TelemetryLogger{
		logger: logger,
	}
}

var _ Telemetry = (*TelemetryLogger)(nil)

func (t *TelemetryLogger) SetGauge(name string, value int64, tags map[string]string) {
	t.logger.Printf("set gauge %s to: %d tags:%v", name, value, tags)
}

func (t *TelemetryLogger) AddGauge(name string, value int64, tags map[string]string) {
	t.logger.Printf("add gauge %s: %d tags:%v", name, value, tags)
}

func (t *TelemetryLogger) SubGauge(name string, value int64, tags map[string]string) {
	t.logger.Printf("sub gauge %s: %d tags:%v", name, value, tags)
}

func (t *TelemetryLogger) AddCounter(name string, value int64, tags map[string]string) {
	t.logger.Printf("counter %s: %d tags:%v", name, value, tags)
}
