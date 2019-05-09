package request

import (
	"strconv"
	"sync"

	"github.com/BiteBit/gorequest"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	Metrics struct {
		Inflight *prometheus.GaugeVec
		Summary  *prometheus.CounterVec
	}
)

var (
	once     sync.Once
	instance *Metrics
)

func getInstance() *Metrics {
	once.Do(func() {
		instance = &Metrics{
			Inflight: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: GlobalOptions.metrics.Namespace,
					Subsystem: GlobalOptions.metrics.Subsystem,
					Name:      "external_request_inflight",
					Help:      "External http request in flight",
				},
				[]string{"method", "url"},
			),
			Summary: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: GlobalOptions.metrics.Namespace,
					Subsystem: GlobalOptions.metrics.Subsystem,
					Name:      "external_request_summary",
					Help:      "External http request summary",
				},
				[]string{"method", "url", "code", "err"},
			),
		}

		prometheus.MustRegister(instance.Inflight, instance.Summary)
	})

	return instance
}

func metricsBefore(agent *gorequest.SuperAgent) {
	if GlobalOptions.Metrics {
		getInstance().Inflight.WithLabelValues(agent.Method, agent.Url).Inc()
	}
}

func metricsAfter(agent *gorequest.SuperAgent, resp gorequest.Response, body []byte, errs []error) {
	if GlobalOptions.Metrics {
		getInstance().Inflight.WithLabelValues(agent.Method, agent.Url).Dec()
		getInstance().Summary.WithLabelValues(agent.Method, agent.Url, strconv.Itoa(resp.StatusCode), formatErrs(errs)).Inc()
	}

}

func formatErrs(errs []error) string {
	errStr := ""
	for _, err := range errs {
		errStr += err.Error() + ";"
	}

	return errStr
}
