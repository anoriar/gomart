package getwithdrawals

import (
	"encoding/json"
	"github.com/anoriar/gophermart/internal/gophermart/balance/handlers/getwithdrawals/internal/factory"
	"github.com/anoriar/gophermart/internal/gophermart/balance/services/withdraw"
	"github.com/anoriar/gophermart/internal/gophermart/shared/context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	reqCtx := req.Context()
	span, reqCtx := opentracing.StartSpanFromContext(reqCtx, "GetWithdrawalsHandler::GetWithdrawals")
	defer span.Finish()

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
		span.LogFields(log.Error(err))
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
