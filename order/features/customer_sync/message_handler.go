package customer_sync

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	customerapi "github.com/pmorelli92/demo-rabbitmq-streams/customer/api"
	"github.com/pmorelli92/demo-rabbitmq-streams/order/metrics"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
)

func HandleMessage(db *sqlx.DB, customerChan chan<- CustomerEvent, message *amqp.Message) {
	var event customerapi.CustomerEvent
	if err := json.Unmarshal(message.Data[0], &event); err != nil {
		log.Printf("Failed to unmarshal customer event: %v", err)
		metrics.CustomerEventErrors.Inc()
		return
	}

	customerEvent := CustomerEvent{
		EventID:    event.EventID,
		EventType:  event.EventType,
		CustomerID: event.CustomerID,
		Timestamp:  event.Timestamp,
		Data:       event.Data,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h := &handler{db: db}
	err := h.processCustomerEvent(ctx, customerEvent)
	if err != nil {
		log.Printf("Failed to process customer event: %v", err)
		metrics.CustomerEventErrors.Inc()
		return
	}
}
