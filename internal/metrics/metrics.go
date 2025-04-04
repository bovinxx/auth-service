package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "auth_space"
	appName   = "auth"
)

type Metrics struct {
	requestCounter        prometheus.Counter
	responseCounter       *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec
}

var metrics *Metrics

func Init(_ context.Context) error {
	metrics = &Metrics{
		requestCounter: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_requests_total",
				Help:      "Amount of requests",
			},
		),
		responseCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_response_total",
				Help:      "Amount of response",
			},
			[]string{"status", "method"},
		),
		histogramResponseTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      appName + "_histogram_response_time_seconds",
				Help:      "Time of response",
				Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
			},
			[]string{"status"},
		),
	}

	prometheus.MustRegister(metrics.requestCounter)
	prometheus.MustRegister(metrics.responseCounter)
	prometheus.MustRegister(metrics.histogramResponseTime)

	return nil
}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}

func IncResponseCounter(status string, method string) {
	metrics.responseCounter.WithLabelValues(status, method).Inc()
}

func HistogramResponseTimeObserve(status string, time float64) {
	metrics.histogramResponseTime.WithLabelValues(status).Observe(time)
}
