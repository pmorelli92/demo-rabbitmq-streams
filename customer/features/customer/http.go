package customer

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type httpHandler struct {
	logger  *slog.Logger
	handler handler
}

func (h httpHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var rq createCustomerRq
	if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
		http.Error(w, "invalid rq body", http.StatusBadRequest)
		return
	}

	rs, err := h.handler.createCustomer(ctx, rq)
	switch err {
	case nil:
		h.logger.InfoContext(ctx, "customer created")
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(rs)
	default:
		h.logger.ErrorContext(ctx, err.Error())
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	}
}

func (h httpHandler) updateCustomerAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerID := r.PathValue("id")

	var rq updateAddressRq
	if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
		http.Error(w, "invalid rq body", http.StatusBadRequest)
		return
	}

	rs, err := h.handler.updateCustomerAddress(ctx, customerID, rq)
	switch err {
	case nil:
		h.logger.InfoContext(ctx, "customer address updated")
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(rs)
	default:
		h.logger.ErrorContext(ctx, err.Error())
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	}
}
