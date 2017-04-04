package metrics

import "github.com/prometheus/client_golang/prometheus"

// Metric string bases for Prometheus Metric names
const (
	MetricAPIGeneralPing      = "probe_general_ping"
	MetricAPIGeneralLogin     = "probe_general_login_duration_seconds"
	MetricAPIChannelGetByName = "probe_channel_getbyname_duration_seconds"
	MetricAPIChannelJoin      = "probe_channel_join_duration_seconds"
	MetricAPIPostCreate       = "probe_post_create"
	MetricProbeBroadcast      = "probe_broadcast_post"
)

//TODO: Let stan now that MetricBrokeBroadcast "probe_broadcast_post_recieve_seconds" -> "probe_broadcast_post_duration_seconds"

// ReponseMetrics holds the response metrics in a map for easy lookup, should match error metrics 1:1
var ResponseMetrics = map[string]prometheus.Gauge{
	// TODO: Figure out if it should be histogram or gurae
	MetricAPIGeneralPing: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: appendResponseSuffix(MetricAPIGeneralPing),
		Help: "Response time of general ping",
	}),
	MetricAPIGeneralLogin: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: appendResponseSuffix(MetricAPIGeneralLogin),
		Help: "Response time of general login",
	}),
	MetricAPIChannelGetByName: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: appendResponseSuffix(MetricAPIChannelGetByName),
		Help: "Response time of channel get by name",
	}),
	MetricAPIChannelJoin: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: MetricAPIChannelJoin,
		Help: "Response time of channel join",
	}),
	MetricAPIPostCreate: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: appendResponseSuffix(MetricAPIPostCreate),
		Help: "Response time of post create",
	}),
	MetricProbeBroadcast: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: appendResponseSuffix(MetricProbeBroadcast),
		Help: "Time from post create to reception by different user",
	}),
}

// ErrorMetrics holds the error metrics for easy lookup, should match response metrics 1:1
var ErrorMetrics = map[string]prometheus.Counter{
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

func appendResponseSuffix(s string) string {
	return s + "_duration_second"
}

func appendErrorSuffix(s string) string {
	return s + "_errors"
}
