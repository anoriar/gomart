package withdraw

import (
	"encoding/json"
	"errors"
	withdrawRequestPkg "github.com/anoriar/gophermart/internal/gophermart/balance/dto/requests"
	errors3 "github.com/anoriar/gophermart/internal/gophermart/balance/errors"
	"github.com/anoriar/gophermart/internal/gophermart/balance/services/withdraw"
	errors2 "github.com/anoriar/gophermart/internal/gophermart/order/errors"
	"github.com/anoriar/gophermart/internal/gophermart/shared/context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	reqCtx := req.Context()
	span, reqCtx := opentracing.StartSpanFromContext(reqCtx, "WithdrawHandler::Withdraw")
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

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		span.LogFields(log.Error(err))
		return
	}
	var withdrawDto withdrawRequestPkg.WithdrawDto
	err = json.Unmarshal(reqBody, &withdrawDto)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		span.LogFields(log.Error(err))
		return
	}

	err = handler.withdrawService.Withdraw(reqCtx, userID, withdrawDto)
	if err != nil {
		switch {
		case errors.Is(err, errors2.ErrOrderNumberNotValid):
			http.Error(w, "not valid order number", http.StatusUnprocessableEntity)
			return
		case errors.Is(err, errors3.ErrInsufficientFunds):
			http.Error(w, "insufficient funds", http.StatusPaymentRequired)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			span.LogFields(log.Error(err))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
