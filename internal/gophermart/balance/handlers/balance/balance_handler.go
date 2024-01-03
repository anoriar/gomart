package balance

import (
	"encoding/json"
	"github.com/anoriar/gophermart/internal/gophermart/balance/services/balance"
	"github.com/anoriar/gophermart/internal/gophermart/shared/context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

type BalanceHandler struct {
	balanceService balance.BalanceServiceInterface
}

func NewBalanceHandler(orderService balance.BalanceServiceInterface) *BalanceHandler {
	return &BalanceHandler{balanceService: orderService}
}

func (handler *BalanceHandler) GetUserBalance(w http.ResponseWriter, req *http.Request) {
	reqCtx := req.Context()
	span, reqCtx := opentracing.StartSpanFromContext(reqCtx, "BalanceHandler::GetUserBalance")
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

	userBalance, err := handler.balanceService.GetUserBalance(reqCtx, userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		span.LogFields(log.Error(err))
		return
	}

	responseBody, err := json.Marshal(userBalance)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		span.LogFields(log.Error(err))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBody)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		span.LogFields(log.Error(err))
		return
	}
}
