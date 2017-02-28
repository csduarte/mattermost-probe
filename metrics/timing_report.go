package metrics

// TimingReport holds metric name and request duration
type TimingReport struct {
	MetricName      string
	Path            string
	DurationSeconds float64
}

// TimingChannel reprsents an incoming report channel
type TimingChannel chan TimingReport
