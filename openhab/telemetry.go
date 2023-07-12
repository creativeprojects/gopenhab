package openhab

const (
	MetricItemName         = "item_name"
	MetricRuleID           = "rule_id"
	MetricItemCacheHit     = "item.cache_hit"
	MetricItemLoad         = "item.load"
	MetricItemLoadState    = "item.load_state"
	MetricItemSetState     = "item.set_state"
	MetricItemNotFound     = "item.not_found"
	MetricItemStateUpdated = "item.state_updated"
	MetricItemsCacheSize   = "items.cache_size"
	MetricRuleAdded        = "rule.added"
	MetricRuleDeleted      = "rule.deleted"
	MetricRulesCount       = "rules.count"
)

// Telemetry interface to send metrics. Two metrics are available: Gauge and Counter.
//
// Gauge is a Metric that represents a single numerical value that can arbitrarily go up and down.
// A Gauge is typically used for measured values like temperatures or current memory usage,
// but also "counts" that can go up and down, like the number of running goroutines.
//
// Counter is a Metric that represents a single numerical value that only ever goes up.
// That implies that it cannot be used to count items whose number can also go down,
// e.g. the number of currently running goroutines. Those "counters" are represented by Gauges.
// A Counter is typically used to count requests served, tasks completed, errors occurred, etc.
type Telemetry interface {
	Close()
	// SetGauge sets the value of a gauge. The callback runs inside its own goroutine
	SetGauge(name string, value int64, tags map[string]string)
	// AddGauge adds the value to a gauge. The callback runs inside its own goroutine
	AddGauge(name string, value int64, tags map[string]string)
	// SubGauge subtracts the value from a gauge. The callback runs inside its own goroutine
	SubGauge(name string, value int64, tags map[string]string)
	// AddCounter adds the value to a counter. The callback runs inside its own goroutine
	AddCounter(name string, value int64, tags map[string]string)
}
