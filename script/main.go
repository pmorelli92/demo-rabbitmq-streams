package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	customerBaseURL = "http://localhost:8080"
	orderBaseURL    = "http://localhost:8081"
)

func main() {
	logger := log.Default()

	logger.Println("1. Create customer with unknown as address")
	rs, err := http.Post(
		customerBaseURL+"/customers",
		"application/json",
		bytes.NewBuffer([]byte(`{
			"name": "John",
			"address": "unknown",
			"email: "john@doe.com"
		}`)))
	if err != nil {
		logger.Fatal(err)
	}
	if rs.StatusCode != http.StatusCreated {
		logger.Fatal(rs.StatusCode)
	}
	defer rs.Body.Close()
	rsBytes, _ := io.ReadAll(rs.Body)
	logger.Println("Customer created -> ", string(rsBytes))

	rsMap := map[string]string{}
	err = json.Unmarshal(rsBytes, &rsMap)
	if err != nil {
		logger.Fatal(err)
	}

	/// ----------

	logger.Println("2. Create order with invalid customer")
	rs, err = http.Post(
		orderBaseURL+"/orders",
		"application/json",
		bytes.NewBuffer([]byte(`{
			"customer_id": "foobar"
		}`)))
	if err != nil {
		logger.Fatal(err)
	}
	if rs.StatusCode != http.StatusBadRequest {
		logger.Fatal(rs.StatusCode)
	}
	logger.Println("Order failed with bad request")

	/// ----------

	logger.Println("3. Create order with valid customer but unknown address")
	rs, err = http.Post(
		orderBaseURL+"/orders",
		"application/json",
		bytes.NewBuffer([]byte(`{
			"customer_id": "`+rsMap["customer_id"]+`"
		}`)))
	if rs.StatusCode != http.StatusCreated {
		logger.Fatal(rs.StatusCode)
	}
	defer rs.Body.Close()
	rsBytes, _ = io.ReadAll(rs.Body)
	logger.Println("Order created -> ", string(rsBytes))

	/// ----------

	logger.Println("4. Update customer address")
	rq, _ := http.NewRequest(
		http.MethodPut,
		customerBaseURL+"/customers"+rsMap["customer_id"]+"/address",
		bytes.NewBuffer([]byte(`{
			"address": "fake street 123"
		}`)))

	rs, err = http.DefaultClient.Do(rq)
	if err != nil {
		logger.Fatal(err)
	}
	if rs.StatusCode != http.StatusOK {
		logger.Fatal(rs.StatusCode)
	}
	logger.Println("Customer address has changed")

	/// ----------

	logger.Println("5. Create order with valid customer that now updated address")
	rs, err = http.Post(
		orderBaseURL+"/orders",
		"application/json",
		bytes.NewBuffer([]byte(`{
			"customer_id": "`+rsMap["customer_id"]+`"
		}`)))
	if rs.StatusCode != http.StatusCreated {
		logger.Fatal(rs.StatusCode)
	}
	defer rs.Body.Close()
	rsBytes, _ = io.ReadAll(rs.Body)
	logger.Println("Order created -> ", string(rsBytes))
}
