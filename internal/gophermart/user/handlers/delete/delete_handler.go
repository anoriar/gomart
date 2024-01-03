package delete

import (
	"github.com/anoriar/gophermart/internal/gophermart/shared/context"
	"github.com/anoriar/gophermart/internal/gophermart/user/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

type DeleteHandler struct {
	userRepository repository.UserRepositoryInterface
}

func NewDeleteHandler(userRepository repository.UserRepositoryInterface) *DeleteHandler {
	return &DeleteHandler{userRepository: userRepository}
}

func (handler *DeleteHandler) Delete(w http.ResponseWriter, req *http.Request) {
	reqCtx := req.Context()
	span, reqCtx := opentracing.StartSpanFromContext(reqCtx, "DeleteHandler::Delete")
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

	err := handler.userRepository.DeleteUser(reqCtx, userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		span.LogFields(log.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
