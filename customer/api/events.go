package api

import "time"

const (
	StreamName             = "customer"
	CustomerCreated        = "customer.created"
	CustomerAddressUpdated = "customer.address.updated"
)

type CustomerCreatedEvent struct {
	CustomerID string    `json:"customer_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Address    string    `json:"address"`
	Timestamp  time.Time `json:"timestamp"`
}

type CustomerAddressUpdatedEvent struct {
	CustomerID string    `json:"customer_id"`
	Address    string    `json:"address"`
	Timestamp  time.Time `json:"timestamp"`
}
