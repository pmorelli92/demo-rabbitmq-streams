package customer

type createCustomerRq struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type updateAddressRq struct {
	Address string `json:"address"`
}

type customerRs struct {
	CustomerID string `json:"customerId"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Address    string `json:"address"`
}
