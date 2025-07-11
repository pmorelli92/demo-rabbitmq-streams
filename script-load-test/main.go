package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type createCustomerRq struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

func main() {
	ctx := context.Background()

	// Connect to both databases
	customerDB, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/customer_db")
	if err != nil {
		log.Fatal("Failed to connect to customer_db:", err)
	}
	defer customerDB.Close()

	orderDB, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/order_db")
	if err != nil {
		log.Fatal("Failed to connect to order_db:", err)
	}
	defer orderDB.Close()

	// Truncate both customer tables
	_, err = customerDB.Exec(ctx, "TRUNCATE TABLE customers")
	if err != nil {
		log.Fatal("Failed to truncate customer_db.customers:", err)
	}
	log.Println("Truncated customer_db.customers")

	_, err = orderDB.Exec(ctx, "TRUNCATE TABLE customers")
	if err != nil {
		log.Fatal("Failed to truncate order_db.customers:", err)
	}
	log.Println("Truncated order_db.customers")

	// Create HTTP client with connection reuse
	client := &http.Client{
		Timeout: 1 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			DisableKeepAlives:   false,
		},
	}

	// Create 100k customers with rate limiting
	const totalCustomers = 100000
	const rateLimit = 20 * time.Millisecond

	log.Printf("Starting to create %d customers...", totalCustomers)

	var wg sync.WaitGroup
	sem := make(chan struct{}, 25) // Limit concurrent requests

	start := time.Now()

	for i := 0; i < totalCustomers; i++ {
		wg.Add(1)
		go func(customerNum int) {
			defer wg.Done()

			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			// Create customer data
			customer := createCustomerRq{
				Name:    fmt.Sprintf("Customer %d", customerNum),
				Email:   fmt.Sprintf("customer%d@example.com", customerNum),
				Address: fmt.Sprintf("Address %d", customerNum),
			}

			// Send HTTP POST request
			jsonData, err := json.Marshal(customer)
			if err != nil {
				log.Printf("Failed to marshal customer %d: %v", customerNum, err)
				return
			}

			resp, err := client.Post("http://localhost:8080/customers", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("Failed to create customer %d: %v", customerNum, err)
				return
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				log.Printf("Unexpected status code %d for customer %d", resp.StatusCode, customerNum)
				return
			}

			// Rate limiting
			time.Sleep(rateLimit)

			// Progress logging
			if customerNum != 0 && customerNum%10000 == 0 {
				elapsed := time.Since(start)
				log.Printf("Created %d customers in %v", customerNum, elapsed)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Finished creating %d customers in %v", totalCustomers, elapsed)

	// Wait a bit for eventual consistency
	log.Println("Waiting 1 second for eventual consistency...")
	time.Sleep(1 * time.Second)

	// Count customers in customer_db
	var customerDBCount int
	err = customerDB.QueryRow(ctx, "SELECT COUNT(*) FROM customers").Scan(&customerDBCount)
	if err != nil {
		log.Fatal("Failed to count customers in customer_db:", err)
	}
	log.Printf("Customer DB count: %d", customerDBCount)

	// Count customers in order_db
	var orderDBCount int
	err = orderDB.QueryRow(ctx, "SELECT COUNT(*) FROM customers").Scan(&orderDBCount)
	if err != nil {
		log.Fatal("Failed to count customers in order_db:", err)
	}
	log.Printf("Order DB count: %d", orderDBCount)

	// Verify counts
	if customerDBCount == totalCustomers && orderDBCount == totalCustomers {
		log.Printf("SUCCESS: Both databases have %d customers as expected", totalCustomers)
	} else {
		log.Printf("FAILURE: Expected %d customers in both databases, got customer_db=%d, order_db=%d",
			totalCustomers, customerDBCount, orderDBCount)
	}
}
