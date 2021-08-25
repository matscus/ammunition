package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	RequestCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ammunition_http_requests_count_total",
		Help: "The total number of request events",
	})
	RequestsDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "ammunition_http_requests_duration_ms",
			Help:       "A summary of the handling duration of requests.",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.01, 0.99: 0.01},
		},
		[]string{"method", "path"},
	)
	WorkerCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ammunition_worker_count_total",
			Help: "Total number of workers",
		},
		[]string{"cache"},
	)
	WorkerDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "ammunition_worker_duration_ms",
			Help:       "A summary of the duration of worker.",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.01, 0.99: 0.01},
		},
		[]string{"cache"},
	)
	CacheCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ammunition_cache_count",
			Help: "cache count",
		},
		[]string{"type", "cache"},
	)
	CacheLen = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ammunition_cache_len_total",
			Help: "cache len",
		},
		[]string{"type", "cache"},
	)
	CacheCap = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ammunition_cache_cap_total",
			Help: "cache len",
		},
		[]string{"type", "cache"},
	)
	Uptime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ammunition_uptime",
			Help: "ammunition uptime.",
		},
	)
)

func init() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(RequestsDuration)
	prometheus.MustRegister(WorkerDuration)
	prometheus.MustRegister(WorkerCount)
	prometheus.MustRegister(CacheCount)
	prometheus.MustRegister(CacheCap)
	prometheus.MustRegister(CacheLen)
	prometheus.MustRegister(Uptime)
	log.Info("Register metrics completed")
}
