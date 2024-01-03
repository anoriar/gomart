package getorders

import (
	"encoding/json"
	"github.com/anoriar/gophermart/internal/gophermart/order/handlers/getorders/internal/factory"
	"github.com/anoriar/gophermart/internal/gophermart/order/services"
	"github.com/anoriar/gophermart/internal/gophermart/shared/context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

type GetOrdersHandler struct {
	orderService         services.OrderServiceInterface
	orderResponseFactory *factory.GetOrdersResponseFactory
}

func NewGetOrdersHandler(orderService services.OrderServiceInterface) *GetOrdersHandler {
	return &GetOrdersHandler{orderService: orderService, orderResponseFactory: factory.NewGetOrdersResponseFactory()}
}

func (handler *GetOrdersHandler) GetOrders(w http.ResponseWriter, req *http.Request) {
	reqCtx := req.Context()
	span, reqCtx := opentracing.StartSpanFromContext(reqCtx, "GetOrdersHandler::GetOrders")
	defer span.Finish()

	userID := ""
	userIDCtxParam := reqCtx.Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	if userID == "" {
		http.Error(w, "user unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := handler.orderService.GetUserOrders(reqCtx, userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		span.LogFields(log.Error(err))
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
		span.LogFields(log.Error(err))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBody)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		span.LogFields(log.Error(err))
		return
	}
}
