package api

import "time"

const (
	StreamName             = "customer"
	CustomerCreated        = "customer.created"
	CustomerAddressUpdated = "customer.address.updated"
)

type CustomerEvent struct {
	EventID    string    `json:"eventId"`
	EventType  string    `json:"eventType"`
	Timestamp  time.Time `json:"timestamp"`
	CustomerID string    `json:"customerId"`
	Data       any       `json:"data"`
}

type CustomerCreatedEvent struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type CustomerAddressUpdatedEvent struct {
	Address string `json:"address"`
}
