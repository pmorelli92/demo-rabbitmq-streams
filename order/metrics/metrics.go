package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	OrdersCreated       prometheus.Counter
	CustomersProcessed  prometheus.Counter
	CustomerEventErrors prometheus.Counter
)

func New() {
	OrdersCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_created_total",
		Help: "The total number of orders created",
	})

	CustomersProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "customers_processed_total",
		Help: "The total number of customer events processed",
	})

	CustomerEventErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "customer_event_errors_total",
		Help: "The total number of customer event processing errors",
	})
}
