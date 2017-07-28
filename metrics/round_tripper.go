package metrics

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

// TimedRoundTripper replacement http.RoundTripper
type TimedRoundTripper struct {
	baseRoundTripper http.RoundTripper
	reportChannel    chan Report
	log              *logrus.Logger
}

// NewTimedRoundTripper will create a new TimedRoundTripper
func NewTimedRoundTripper(trc chan Report, log *logrus.Logger) *TimedRoundTripper {
	rt := TimedRoundTripper{
		http.DefaultTransport,
		trc,
		log,
	}

	return &rt
}

// RoundTrip will send off the response time to the report channel
func (trt TimedRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	requestStart := time.Now()
	resp, err := trt.baseRoundTripper.RoundTrip(r)
	requestEnd := time.Now()
	requestDuration := requestEnd.Sub(requestStart).Seconds()

	if err != nil {
		err = errors.Wrap(err, "Http error respose")
		trt.log.Errorf("TimedRoundTripper detected HTTP error: %s", err.Error())
	}

	if resp != nil && resp.StatusCode >= 400 {
		var body string
		if resp.Body != nil {
			data, rErr := ioutil.ReadAll(resp.Body)
			if rErr != nil {
				trt.log.Errorf("TimedRoundTrip failed to read response body of error")
			} else {
				bErr := resp.Body.Close()
				if bErr != nil {
					trt.log.Errorf("TimedRoundTripper failed to close Body")
				}
				body = string(data)
			}
		}
		err = errors.New("Response code greater than 399, forcing error")
		trt.log.Errorf("TimedRoundTripper detected response >= 400. Response Body: \n%q", body)
	}

	if time.Duration(requestDuration) > 10*time.Second {
		err = errors.New("Response Duration >= 10s, forcing error")
	}

	var route string
	if pingRoute := r.Header.Get("PingProbe"); len(pingRoute) > 0 {
		route = pingRoute
	} else {
		route = TokenizePath(r.URL.Path)
	}

	trt.reportChannel <- Report{
		Route:           route,
		DurationSeconds: requestDuration,
		Error:           err,
	}

	return resp, err
}
