package monitoring

import "github.com/prometheus/client_golang/prometheus"

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
	WorkerDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "ammunition_worker_duration_ms",
			Help:       "A summary of the duration of worker.",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.01, 0.99: 0.01},
		},
		[]string{"cache_name"},
	)
	CacheCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ammunition_cache_count",
			Help: "cache count",
		},
	)
	CacheLen = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ammunition_cache_len_total",
			Help: "cache len",
		},
		[]string{"cache_name"},
	)
	Uptime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ammunition_uptime",
			Help: "ammunition uptime.",
		},
	)
)

func init() {
	registry := prometheus.NewRegistry()
	registry.Register(RequestCount)
	registry.Register(RequestsDuration)
	registry.Register(Uptime)
}
