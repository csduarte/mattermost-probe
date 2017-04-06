package metrics

import "github.com/prometheus/client_golang/prometheus"

// Metric string bases for Prometheus Metric names
const (
	MetricAPIGeneralPing      = "probe_general_ping"
	MetricAPIGeneralLogin     = "probe_general_login"
	MetricAPIChannelGetByName = "probe_channel_getbyname"
	MetricAPIChannelJoin      = "probe_channel_join"
	MetricAPIPostCreate       = "probe_post_create"
	MetricProbeBroadcast      = "probe_broadcast_post"
)

//TODO: Let stan now that MetricBrokeBroadcast "probe_broadcast_post_recieve_seconds" -> "probe_broadcast_post_duration_seconds"

// ResponseHistograms holds the response metrics in a map for easy lookup, should match error metrics 1:1
var ResponseHistograms = map[string]prometheus.Histogram{
	MetricAPIGeneralPing: prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: appendHistogramSuffix(MetricAPIGeneralPing),
		Help: "Response time of general ping",
	}),
	MetricAPIGeneralLogin: prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: appendHistogramSuffix(MetricAPIGeneralLogin),
		Help: "Response time of general login",
	}),
	MetricAPIChannelGetByName: prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: appendHistogramSuffix(MetricAPIChannelGetByName),
		Help: "Response time of channel get by name",
	}),
	MetricAPIChannelJoin: prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: appendHistogramSuffix(MetricAPIChannelJoin),
		Help: "Response time of channel join",
	}),
	MetricAPIPostCreate: prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: appendHistogramSuffix(MetricAPIPostCreate),
		Help: "Response time of post create",
	}),
	MetricProbeBroadcast: prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: appendHistogramSuffix(MetricProbeBroadcast),
		Help: "Time from post create to reception by different user",
	}),
}

//ResponseGauges holds response gauages for latency
var ResponseGauges = map[string]prometheus.Gauge{
	MetricAPIGeneralPing: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: appendGaugeSuffix(MetricAPIGeneralPing),
		Help: "Response time of general ping",
	}),
	MetricAPIGeneralLogin: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: appendGaugeSuffix(MetricAPIGeneralLogin),
		Help: "Response time of general login",
	}),
	MetricAPIChannelGetByName: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: appendGaugeSuffix(MetricAPIChannelGetByName),
		Help: "Response time of channel get by name",
	}),
	MetricAPIChannelJoin: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: appendGaugeSuffix(MetricAPIChannelJoin),
		Help: "Response time of channel join",
	}),
	MetricAPIPostCreate: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: appendGaugeSuffix(MetricAPIPostCreate),
		Help: "Response time of post create",
	}),
	MetricProbeBroadcast: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: appendGaugeSuffix(MetricProbeBroadcast),
		Help: "Time from post create to reception by different user",
	}),
}

// ErrorCounters holds the error metrics for easy lookup, should match response metrics 1:1
var ErrorCounters = map[string]prometheus.Counter{
	// TODO: Figure out if it should be histogram or gurae
	MetricAPIGeneralPing: prometheus.NewCounter(prometheus.CounterOpts{
		Name: appendErrorSuffix(MetricAPIGeneralPing),
		Help: "Errors for general ping",
	}),
	MetricAPIGeneralLogin: prometheus.NewCounter(prometheus.CounterOpts{
		Name: appendErrorSuffix(MetricAPIGeneralLogin),
		Help: "Errors for general login",
	}),
	MetricAPIChannelGetByName: prometheus.NewCounter(prometheus.CounterOpts{
		Name: appendErrorSuffix(MetricAPIChannelGetByName),
		Help: "Errors for channel get by name",
	}),
	MetricAPIChannelJoin: prometheus.NewCounter(prometheus.CounterOpts{
		Name: appendErrorSuffix(MetricAPIChannelJoin),
		Help: "Errors for channel join",
	}),
	MetricAPIPostCreate: prometheus.NewCounter(prometheus.CounterOpts{
		Name: appendErrorSuffix(MetricAPIPostCreate),
		Help: "Errors for post create",
	}),
	MetricProbeBroadcast: prometheus.NewCounter(prometheus.CounterOpts{
		Name: appendErrorSuffix(MetricProbeBroadcast),
		Help: "Errors for Time from post create to reception by different user",
	}),
}

func appendGaugeSuffix(s string) string {
	return s + "_duration_seconds"
}

func appendHistogramSuffix(s string) string {
	return s + "duration_histogram_seconds"
}

func appendErrorSuffix(s string) string {
	return s + "_errors"
}
