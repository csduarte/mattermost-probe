package metrics

import (
	"regexp"
	"strings"

	"github.com/mattermost/platform/model"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricAPIGeneralPing      = "probe_general_ping_duration_seconds"
	MetricAPIGeneralLogin     = "probe_general_login_duration_seconds"
	MetricAPIChannelGetByName = "probe_channel_getbyname_duration_seconds"
	MetricAPIChannelJoin      = "probe_channel_join_duration_seconds"
	MetricAPIPostCreate       = "probe_post_create_duration_seconds"
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
}

var metricNamesByCommonPath = map[string]string{
	"/general/ping":                        MetricAPIGeneralPing,
	"/users/login":                         MetricAPIGeneralLogin,
	"/teams/tid/channels/name/cname":       MetricAPIChannelGetByName,
	"/teams/tid/channels/cid/join":         MetricAPIChannelJoin,
	"/teams/tid/channels/cid/posts/create": MetricAPIPostCreate,
}

// rawSubstitutions holds the raw strings that will be compiled at run time
var rawSubstitutions = map[string]string{
	//Channel Routes
	"/channels/cid/":       "/channels/[a-z0-9]{26}/",       //Channel ID
	"/channels/name/cname": "/channels/name/[A-Za-z0-9_-]+", //Get Channel By Name

	//Team Routes
	"/teams/tid/": "/teams/[a-z0-9]{26}/", //Team ID
}

// Subtitutions holds the compiled regex
var Subtitutions = map[string]*regexp.Regexp{}

func init() {
	for k, v := range rawSubstitutions {
		Subtitutions[k] = regexp.MustCompile(v)
	}
}

// CollatePaths organize and clean up path names based off known formats in Subtitutions
func CollatePaths(path string) string {
	result := strings.TrimPrefix(path, model.API_URL_SUFFIX)
	for sub, reg := range Subtitutions {
		result = reg.ReplaceAllString(result, sub)
	}
	return result
}

// LookupMetricNameByPath will search for metricname by path
func LookupMetricNameByPath(path string) (string, bool) {
	cp := CollatePaths(path)
	name, ok := metricNamesByCommonPath[cp]
	return name, ok
}
