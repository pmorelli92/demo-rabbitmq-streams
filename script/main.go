package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	customerBaseURL = "http://localhost:8080"
	orderBaseURL    = "http://localhost:8081"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("1. Create customer with unknown as address")
	rqBody := map[string]string{
		"name":    "John Doe",
		"email":   "john@example.com",
		"address": "unknown",
	}
	b, _ := json.Marshal(rqBody)
	rs, err := http.Post(customerBaseURL+"/customers", "application/json", bytes.NewBuffer(b))
	if err != nil {
		logger.Error("error", "error", err)
		return
	}
	if rs.StatusCode != http.StatusCreated {
		logger.Error("error", "status code", rs.StatusCode)
		return
	}
	defer rs.Body.Close()
	rsBytes, _ := io.ReadAll(rs.Body)
	logger.Info("1. Customer created", "rs", string(rsBytes))

	rsMap := map[string]string{}
	err = json.Unmarshal(rsBytes, &rsMap)
	if err != nil {
		logger.Error("error", "error", err)
		return
	}

	time.Sleep(100 * time.Millisecond)
	/// ----------

	logger.Info("2. Create order with invalid customer")
	rqBody = map[string]string{
		"customer_id": "foo bar",
	}
	b, _ = json.Marshal(rqBody)
	rs, err = http.Post(orderBaseURL+"/orders", "application/json", bytes.NewBuffer(b))
	if err != nil {
		logger.Error("error", "error", err)
		return
	}
	if rs.StatusCode != http.StatusBadRequest {
		logger.Error("error", "status code", rs.StatusCode)
		return
	}
	logger.Info("2. Order failed with bad request")

	time.Sleep(100 * time.Millisecond)
	/// ----------

	logger.Info("3. Create order with valid customer but unknown address")
	rqBody = map[string]string{
		"customer_id": rsMap["customer_id"],
	}
	b, _ = json.Marshal(rqBody)
	rs, err = http.Post(orderBaseURL+"/orders", "application/json", bytes.NewBuffer(b))
	if err != nil {
		logger.Error("error", "error", err)
		return
	}
	if rs.StatusCode != http.StatusCreated {
		logger.Error("error", "status code", rs.StatusCode)
	}
	defer rs.Body.Close()
	rsBytes, _ = io.ReadAll(rs.Body)
	logger.Info("3. Order created", "rs", string(rsBytes))

	time.Sleep(100 * time.Millisecond)
	/// ----------

	logger.Info("4. Update customer address")
	rqBody = map[string]string{
		"address": "fake street 123",
	}
	b, _ = json.Marshal(rqBody)
	rq, _ := http.NewRequest(http.MethodPut, customerBaseURL+"/customers/"+rsMap["customer_id"]+"/address", bytes.NewBuffer(b))

	rs, err = http.DefaultClient.Do(rq)
	if err != nil {
		logger.Error("error", "error", err)
		return
	}
	if rs.StatusCode != http.StatusCreated {
		logger.Error("error", "status code", rs.StatusCode)
	}
	logger.Info("4. Customer address has changed")

	time.Sleep(100 * time.Millisecond)
	/// ----------

	logger.Info("5. Create order with valid customer that now updated address")
	rqBody = map[string]string{
		"customer_id": rsMap["customer_id"],
	}
	b, _ = json.Marshal(rqBody)
	rs, err = http.Post(orderBaseURL+"/orders", "application/json", bytes.NewBuffer(b))
	if err != nil {
		logger.Error("error", "error", err)
		return
	}
	if rs.StatusCode != http.StatusCreated {
		logger.Error("error", "status code", rs.StatusCode)
	}
	defer rs.Body.Close()
	rsBytes, _ = io.ReadAll(rs.Body)
	logger.Info("5. Order created", "rs", string(rsBytes))
}
