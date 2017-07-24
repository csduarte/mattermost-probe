package metrics

import (
	"errors"
	"net/http"
	"time"
)

// TimedRoundTripper replacement http.RoundTripper
type TimedRoundTripper struct {
	baseRoundTripper http.RoundTripper
	reportChannel    chan TimingReport
}

// NewTimedRoundTripper will create a new TimedRoundTripper
func NewTimedRoundTripper(trc chan TimingReport) *TimedRoundTripper {
	rt := TimedRoundTripper{
		http.DefaultTransport,
		trc,
	}

	return &rt
}

// RoundTrip will send off the response time to the report channel
func (trt TimedRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	requestStart := time.Now()
	resp, err := trt.baseRoundTripper.RoundTrip(r)
	requestEnd := time.Now()
	requestDuration := requestEnd.Sub(requestStart).Seconds()

	if err != nil || resp.StatusCode >= 400 {
		err = errors.New("Response Code >= 400, forcing error")
	}
	if time.Duration(requestDuration) > 10*time.Second {
		err = errors.New("Response Duration >= 10s, forcing error")
	}

	trt.reportChannel <- TimingReport{
		"",
		r.URL.Path,
		requestDuration,
		err,
	}

	return resp, err
}
