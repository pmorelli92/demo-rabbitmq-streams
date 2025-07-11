package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CustomersCreated   prometheus.Counter
	CustomersUpdated   prometheus.Counter
	EventsPublished    prometheus.Counter
	EventPublishErrors prometheus.Counter
)

func New() {
	CustomersCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "customers_created_total",
		Help: "The total number of customers created",
	})

	CustomersUpdated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "customers_updated_total",
		Help: "The total number of customers updated",
	})

	EventsPublished = promauto.NewCounter(prometheus.CounterOpts{
		Name: "events_published_total",
		Help: "The total number of events published",
	})

	EventPublishErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "event_publish_errors_total",
		Help: "The total number of event publish errors",
	})
}
