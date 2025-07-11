package order

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type httpHandler struct {
	logger  *slog.Logger
	handler handler
}

func (h httpHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var rq createOrderRq
	if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
		http.Error(w, "invalid rq body", http.StatusBadRequest)
		return
	}

	rs, err := h.handler.createOrder(ctx, rq)
	switch err {
	case nil:
		h.logger.InfoContext(ctx, "order created")
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(rs)
	default:
		h.logger.ErrorContext(ctx, err.Error())
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	}
}

func (h httpHandler) getOrderByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderID := r.PathValue("id")

	rs, err := h.handler.getOrder(ctx, orderID)
	switch err {
	case nil:
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(rs)
	default:
		h.logger.ErrorContext(ctx, err.Error())
		http.Error(w, "unexpected error", http.StatusInternalServerError)
	}
}
