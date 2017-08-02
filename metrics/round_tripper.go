package metrics

import (
	"fmt"
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
	origin           string
}

// NewTimedRoundTripper will create a new TimedRoundTripper
func NewTimedRoundTripper(trc chan Report, log *logrus.Logger, origin string) *TimedRoundTripper {
	rt := TimedRoundTripper{
		http.DefaultTransport,
		trc,
		log,
		origin,
	}

	return &rt
}

// RoundTrip will send off the response time to the report channel
func (trt TimedRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if len(trt.origin) > 0 {
		r.Header.Add("Origin", trt.origin)
	}
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
		errMsg := fmt.Sprintf("TimedRoundTripper detected response >= 400. Response Code: %d Response Body: \n%q", resp.StatusCode, body)
		err = errors.New(errMsg)
		trt.log.Errorf(errMsg)
	}

	if time.Duration(requestDuration) > 10*time.Second {
		err = errors.New("Response Duration >= 10s, forcing error")
	}

	var route string
	if pingRoute := r.Header.Get("ProbeRouteOverride"); len(pingRoute) > 0 {
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
