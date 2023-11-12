package load_order

import (
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/context"
	"github.com/anoriar/gophermart/internal/gophermart/domain_errors"
	"github.com/anoriar/gophermart/internal/gophermart/services/order"
	"io"
	"net/http"
)

type LoadOrderHandler struct {
	orderService order.OrderServiceInterface
}

func NewLoadOrderHandler(orderService order.OrderServiceInterface) *LoadOrderHandler {
	return &LoadOrderHandler{orderService: orderService}
}

func (handler *LoadOrderHandler) LoadOrder(w http.ResponseWriter, req *http.Request) {
	userID := ""
	userIDCtxParam := req.Context().Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	if userID == "" {
		http.Error(w, "user unauthorized", http.StatusUnauthorized)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "text/plain" {
		http.Error(w, "not valid request format", http.StatusBadRequest)
		return
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	orderID := string(reqBody)

	err = handler.orderService.LoadOrder(req.Context(), orderID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain_errors.ErrOrderNumberNotValid):
			http.Error(w, "not valid order number", http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domain_errors.ErrOrderAlreadyLoaded):
			w.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, domain_errors.ErrOrderLoadedByAnotherUser):
			http.Error(w, "order already loaded by another user", http.StatusConflict)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
	return
}
