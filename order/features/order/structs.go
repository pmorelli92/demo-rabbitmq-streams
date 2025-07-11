package order

type createOrderRq struct {
	Name       string `json:"name"`
	CustomerID string `json:"customer_id"`
}

type createOrderRs struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type orderRs struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
}
