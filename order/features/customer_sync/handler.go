package customer_sync

import (
	"context"
	"fmt"
	"log"

	customerapi "github.com/pmorelli92/demo-rabbitmq-streams/customer/api"
)

func (h *handler) processCustomerEvent(ctx context.Context, event CustomerEvent) error {
	switch event.EventType {
	case customerapi.CustomerCreated:
		fmt.Println("Processing CustomerCreated event")
		return h.handleCustomerCreated(ctx, event)
	case customerapi.CustomerAddressUpdated:
		fmt.Println("Processing CustomerAddressUpdated event")
		return h.handleCustomerAddressUpdated(ctx, event)
	default:
		log.Printf("Unknown customer event type: %s", event.EventType)
		return nil
	}
}

func (h *handler) handleCustomerCreated(ctx context.Context, event CustomerEvent) error {
	data := event.Data.(map[string]interface{})
	customer := Customer{
		ID:        event.CustomerID,
		Name:      data["name"].(string),
		Email:     data["email"].(string),
		Address:   data["address"].(string),
		CreatedAt: event.Timestamp,
		UpdatedAt: event.Timestamp,
	}

	return h.upsertCustomer(ctx, customer)
}

func (h *handler) handleCustomerAddressUpdated(ctx context.Context, event CustomerEvent) error {
	data := event.Data.(map[string]interface{})
	newAddress := data["new_address"].(string)

	return h.updateCustomerAddress(ctx, event.CustomerID, newAddress)
}
