package metrics

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

// Server respresents the Prometheus metrics and incoming channel
type Server struct {
	ReportChannel chan Report
	Log           *logrus.Logger
	MetricsLog    *logrus.Logger
}

// NewServer returns a new metric server that is ready processing timing reports
func NewServer(log, mlog *logrus.Logger) *Server {
	tr := make(chan Report)

	s := &Server{
		ReportChannel: tr,
		Log:           log,
		MetricsLog:    mlog,
	}

	go MonitorTimingReports(tr, log, mlog)

	return s
}

// Listen starts the prometheus server
func (s *Server) Listen(address string, port int) {
	if port == 0 {
		port = 8067
	}
	if address == "" {
		address = "0.0.0.0"
	}

	RegisterMetrics()

	serverAddr := fmt.Sprintf("%s:%d", address, port)
	if s.Log != nil {
		s.Log.Info("Metrics Server Started on ", serverAddr)
	}

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(serverAddr, nil)

}
