package metrics

import "github.com/prometheus/client_golang/prometheus"

// Metric string bases for Prometheus Metric names
const (
	MetricRequestDurationHistogram = "probe_request_duration_histogram_seconds"
	MetricRequestDuration          = "probe_request_duration_seconds"
	MetricErrorCount               = "probe_error_total"
)

// ResponseHistogram holds the response metrics for API calls and pings
var ResponseHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: MetricRequestDurationHistogram,
	Help: "Response Histogram for API and Pings",
},
	[]string{"route"},
)

// ResponseGauge holds the response time for API calls and pings
var ResponseGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: MetricRequestDuration,
	Help: "Response time for API and Pings",
},
	[]string{"route"},
)

// ErrorCounter represents any 400+ response errors from the API
var ErrorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: MetricErrorCount,
	Help: "Errors counts for API and Pings",
},
	[]string{"route"},
)

// MetricsRegistry holds prometheus collectors for registration.
var MetricsRegistry = []prometheus.Collector{
	ResponseHistogram,
	ResponseGauge,
	ErrorCounter,
}

// RegisterMetrics will ensure all metrics are registered with default prometheus handler
func RegisterMetrics() {
	for _, m := range MetricsRegistry {
		prometheus.MustRegister(m)
	}
}
