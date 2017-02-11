package metrics

import (
	"fmt"
	"net/http"
)

// Server respresents the Prometheus metrics and incoming channel
type Server struct {
	ReportChannel TimingChannel
}

// NewServer returns a new metric server that is ready
func NewServer() *Server {
	tr := make(chan TimingReport)
	s := &Server{
		tr,
	}
	go s.MonitorTimingReports()
	return s
}

// Listen starts the prometheus server
func (s *Server) Listen() {
	// http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":8080", nil)
}

//MonitorTimingReports loops through report channel andpasses the reports to HandleReport
func (s *Server) MonitorTimingReports() {
	for {
		select {
		case r := <-s.ReportChannel:
			s.HandleReport(r)
		}
	}
}

// HandleReport will process incoming timing reports
func (s *Server) HandleReport(r TimingReport) {
	mn, ok := LookupMetricNameByPath(r.Path)
	if ok {
		fmt.Printf("%v, %v", mn, r.RequestDuration)
	} else {
		fmt.Printf("WARNING: Failed to find metric by path %v\n", r.Path)
	}
}
