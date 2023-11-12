package get_orders

import (
	"encoding/json"
	"github.com/anoriar/gophermart/internal/gophermart/context"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/get_orders/internal/factory"
	"github.com/anoriar/gophermart/internal/gophermart/services/order"
	"net/http"
)

type GetOrdersHandler struct {
	orderService         order.OrderServiceInterface
	orderResponseFactory *factory.GetOrdersResponseFactory
}

func NewGetOrdersHandler(orderService order.OrderServiceInterface) *GetOrdersHandler {
	return &GetOrdersHandler{orderService: orderService, orderResponseFactory: factory.NewGetOrdersResponseFactory()}
}

func (handler *GetOrdersHandler) GetOrders(w http.ResponseWriter, req *http.Request) {
	userID := ""
	userIDCtxParam := req.Context().Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	if userID == "" {
		http.Error(w, "user unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := handler.orderService.GetUserOrders(req.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response := handler.orderResponseFactory.CreateOrdersResponse(orders)

	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBody)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	return
}
