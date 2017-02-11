package metrics

import "time"

// TimingReport holds metric name and request duration
type TimingReport struct {
	Path            string
	RequestDuration time.Duration
}

// TimingChannel reprsents an incoming report channel
type TimingChannel chan TimingReport
