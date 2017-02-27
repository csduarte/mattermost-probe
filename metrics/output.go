package metrics

import "go.uber.org/zap"

// NewMetricOutput creates a zap logger suitable for timing reports
func NewMetricOutput(output string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = append(cfg.OutputPaths, output)
	cfg.DisableCaller = true
	return cfg.Build()
}
