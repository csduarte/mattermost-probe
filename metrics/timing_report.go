package metrics

import "go.uber.org/zap/zapcore"

// TimingReport holds metric name and request duration
type TimingReport struct {
	MetricName      string
	Path            string
	DurationSeconds float64
}

// TimingChannel reprsents an incoming report channel
type TimingChannel chan TimingReport

// MarshalLogObject conforms to zapcore encoder
func (tr TimingReport) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", tr.MetricName)
	enc.AddFloat64("duration", tr.DurationSeconds)
	return nil
}
