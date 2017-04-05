package metrics

// TimingReport holds metric name and request duration
type TimingReport struct {
	MetricName      string
	Path            string
	DurationSeconds float64
	Error           error
}

// TimingChannel reprsents an incoming report channel
type TimingChannel chan TimingReport

// EnsureName will fetch the name by path if possible, otherwise return empty string and not ok.
func (tr *TimingReport) EnsureName() (string, bool) {
	if len(tr.MetricName) > 0 {
		return tr.MetricName, true
	}
	if len(tr.Path) == 0 {
		return "", false
	}
	mn, ok := LookupMetricNameByPath(tr.Path)
	if !ok {
		return "", false
	}
	tr.MetricName = mn
	return mn, true
}
