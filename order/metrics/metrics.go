package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "mg"
	subsystem = "order"
)

var (
	OrdersCreated         prometheus.Counter
	CustomersProcessed    prometheus.Counter
	CustomerEventErrors   prometheus.Counter
	HTTPRequests          *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	StreamMessagesProcessed prometheus.Counter
	StreamOffsetUpdates   prometheus.Counter
)

func New() {
	OrdersCreated = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "orders_created_total",
		Help:      "The total number of orders created",
	})

	CustomersProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "customers_processed_total",
		Help:      "The total number of customer events processed",
	})

	CustomerEventErrors = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "customer_event_errors_total",
		Help:      "The total number of customer event processing errors",
	})

	HTTPRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "http_requests_total",
		Help:      "The total number of HTTP requests",
	}, []string{"method", "endpoint", "status"})

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "http_request_duration_seconds",
		Help:      "HTTP request duration in seconds",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "endpoint"})

	StreamMessagesProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "stream_messages_processed_total",
		Help:      "The total number of stream messages processed",
	})

	StreamOffsetUpdates = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "stream_offset_updates_total",
		Help:      "The total number of stream offset updates",
	})
}