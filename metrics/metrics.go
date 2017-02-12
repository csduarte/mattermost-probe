package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	MetricAPIGeneralPing      = "probe_general_ping_duration_seconds"
	MetricAPIGeneralLogin     = "probe_general_login_duration_seconds"
	MetricAPIChannelGetByName = "probe_channel_getbyname_duration_seconds"
	MetricAPIChannelJoin      = "probe_channel_join_duration_seconds"
	MetricAPIPostCreate       = "probe_post_create_duration_seconds"
	MetricProbeBroadcast      = "probe_broadcast_post_received_seconds"
)

// Metrics holds the prometheus metrics in a map for easy lookup
var Metrics = map[string]prometheus.Gauge{
	MetricAPIGeneralPing: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: MetricAPIGeneralPing,
		Help: "Response time of general ping",
	}),
	MetricAPIGeneralLogin: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: MetricAPIGeneralLogin,
		Help: "Response time of general login",
	}),
	MetricAPIChannelGetByName: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: MetricAPIChannelGetByName,
		Help: "Response time of channel get by name",
	}),
	MetricAPIChannelJoin: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: MetricAPIChannelJoin,
		Help: "Response time of channel join",
	}),
	MetricAPIPostCreate: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: MetricAPIPostCreate,
		Help: "Response time of post create",
	}),
	MetricProbeBroadcast: prometheus.NewGauge(prometheus.GaugeOpts{
		Name: MetricProbeBroadcast,
		Help: "Time from post create to reception by different user",
	}),
}
