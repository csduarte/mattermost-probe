package metrics

import (
	"net/http"
	"time"
)

// TimedRoundTripper replacement http.RoundTripper
type TimedRoundTripper struct {
	baseRoundTripper http.RoundTripper
	reportChannel    TimingChannel
}

// NewTimedRoundTripper will create a new TimedRoundTripper
func NewTimedRoundTripper(reportChanel chan TimingReport) *TimedRoundTripper {
	rt := &TimedRoundTripper{
		http.DefaultTransport,
		reportChanel,
	}

	return rt
}

// RoundTrip will send off the response time to the report channel
func (trt *TimedRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	requestStart := time.Now()
	resp, err := trt.baseRoundTripper.RoundTrip(r)
	requestEnd := time.Now()
	requestDuration := requestEnd.Sub(requestStart).Seconds()

	if err != nil || resp.StatusCode >= 400 {
		requestDuration = 0
	}

	trt.reportChannel <- TimingReport{
		"",
		r.URL.Path,
		requestDuration,
	}

	return resp, err
}
