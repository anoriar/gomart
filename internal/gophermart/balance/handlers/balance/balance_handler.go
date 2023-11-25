package balance

import (
	"encoding/json"
	"github.com/anoriar/gophermart/internal/gophermart/balance/services/balance"
	"github.com/anoriar/gophermart/internal/gophermart/shared/context"
	"net/http"
)

type BalanceHandler struct {
	balanceService balance.BalanceServiceInterface
}

func NewBalanceHandler(orderService balance.BalanceServiceInterface) *BalanceHandler {
	return &BalanceHandler{balanceService: orderService}
}

func (handler *BalanceHandler) GetUserBalance(w http.ResponseWriter, req *http.Request) {
	userID := ""
	userIDCtxParam := req.Context().Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	if userID == "" {
		http.Error(w, "user unauthorized", http.StatusUnauthorized)
		return
	}

	userBalance, err := handler.balanceService.GetUserBalance(req.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(userBalance)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBody)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
