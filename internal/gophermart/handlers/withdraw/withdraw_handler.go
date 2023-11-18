package withdraw

import (
	"encoding/json"
	"errors"
	"github.com/anoriar/gophermart/internal/gophermart/context"
	"github.com/anoriar/gophermart/internal/gophermart/domainerrors"
	withdrawRequestPkg "github.com/anoriar/gophermart/internal/gophermart/dto/requests/withdraw"
	"github.com/anoriar/gophermart/internal/gophermart/services/withdraw"
	"io"
	"net/http"
)

type WithdrawHandler struct {
	withdrawService withdraw.WithdrawServiceInterface
}

func NewWithdrawHandler(withdrawService withdraw.WithdrawServiceInterface) *WithdrawHandler {
	return &WithdrawHandler{withdrawService: withdrawService}
}

func (handler *WithdrawHandler) Withdraw(w http.ResponseWriter, req *http.Request) {
	userID := ""
	userIDCtxParam := req.Context().Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	if userID == "" {
		http.Error(w, "user unauthorized", http.StatusUnauthorized)
		return
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var withdrawDto withdrawRequestPkg.WithdrawDto
	err = json.Unmarshal(reqBody, &withdrawDto)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = handler.withdrawService.Withdraw(req.Context(), userID, withdrawDto)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrOrderNumberNotValid):
			http.Error(w, "not valid order number", http.StatusUnprocessableEntity)
			return
		case errors.Is(err, domainerrors.ErrInsufficientFunds):
			http.Error(w, "insufficient funds", http.StatusPaymentRequired)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
