package metrics

import "github.com/Sirupsen/logrus"

// Report holds metric name and request duration
type Report struct {
	Route           string
	DurationSeconds float64
	Error           error
}

//MonitorTimingReports loops through report channel andpasses the reports to HandleReport
func MonitorTimingReports(rc chan Report, log, mlog *logrus.Logger) {
	for {
		select {
		case r := <-rc:
			r.Process(log, mlog)
		}
	}
}

// Process will
func (r Report) Process(log, mlog *logrus.Logger) {
	if log != nil {
		log.WithFields(logrus.Fields{
			"Route":    r.Route,
			"Duration": r.DurationSeconds,
			"Error":    r.Error,
		}).Info("Incoming Report")
	}

	if r.Error != nil {
		ErrorCounter.WithLabelValues(r.Route).Inc()
		if mlog != nil {
			mlog.WithFields(logrus.Fields{
				"Route": r.Route,
				"Err":   r.Error,
			}).Debug("metric error")
		}
		return
	}

	if len(r.Route) == 0 {
		if log != nil {
			log.Warn("Report came in with no Route set, will not log")
		}
		return
	}

	ResponseGauge.WithLabelValues(r.Route).Set(r.DurationSeconds)
	ResponseHistogram.WithLabelValues(r.Route).Observe(r.DurationSeconds)

	if mlog != nil {
		mlog.WithFields(logrus.Fields{
			"Route":    r.Route,
			"Duration": r.DurationSeconds,
		}).Info("metric")
	}
}
