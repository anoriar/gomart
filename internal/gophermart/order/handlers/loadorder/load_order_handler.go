package loadorder

import (
	"errors"
	errors2 "github.com/anoriar/gophermart/internal/gophermart/order/errors"
	"github.com/anoriar/gophermart/internal/gophermart/order/services"
	"github.com/anoriar/gophermart/internal/gophermart/shared/context"
	"io"
	"net/http"
)

type LoadOrderHandler struct {
	orderService services.OrderServiceInterface
}

func NewLoadOrderHandler(orderService services.OrderServiceInterface) *LoadOrderHandler {
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
		case errors.Is(err, errors2.ErrOrderNumberNotValid):
			http.Error(w, "not valid order number", http.StatusUnprocessableEntity)
			return
		case errors.Is(err, errors2.ErrOrderAlreadyLoaded):
			w.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, errors2.ErrOrderLoadedByAnotherUser):
			http.Error(w, "order already loaded by another user", http.StatusConflict)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
