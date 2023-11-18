package getwithdrawals

import (
	"encoding/json"
	"github.com/anoriar/gophermart/internal/gophermart/context"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/getwithdrawals/internal/factory"
	"github.com/anoriar/gophermart/internal/gophermart/services/withdraw"
	"net/http"
)

type GetWithdrawalsHandler struct {
	withdrawService           withdraw.WithdrawServiceInterface
	withdrawalResponseFactory *factory.GetWithdrawalsResponseFactory
}

func NewGetWithdrawalsHandler(withdrawService withdraw.WithdrawServiceInterface) *GetWithdrawalsHandler {
	return &GetWithdrawalsHandler{withdrawService: withdrawService, withdrawalResponseFactory: factory.NewGetWithdrawalsResponseFactory()}
}

func (handler *GetWithdrawalsHandler) GetWithdrawals(w http.ResponseWriter, req *http.Request) {
	userID := ""
	userIDCtxParam := req.Context().Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	if userID == "" {
		http.Error(w, "user unauthorized", http.StatusUnauthorized)
		return
	}

	withdrawals, err := handler.withdrawService.GetWithdrawalsByUserID(req.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response := handler.withdrawalResponseFactory.CreateWithdrawalsResponse(withdrawals)

	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBody)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
}
