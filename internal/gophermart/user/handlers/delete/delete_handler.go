package delete

import (
	"github.com/anoriar/gophermart/internal/gophermart/shared/context"
	"github.com/anoriar/gophermart/internal/gophermart/user/repository"
	"net/http"
)

type DeleteHandler struct {
	userRepository repository.UserRepositoryInterface
}

func NewDeleteHandler(userRepository repository.UserRepositoryInterface) *DeleteHandler {
	return &DeleteHandler{userRepository: userRepository}
}

func (handler *DeleteHandler) Delete(w http.ResponseWriter, req *http.Request) {

	userID := ""
	userIDCtxParam := req.Context().Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}

	if userID == "" {
		http.Error(w, "user unauthorized", http.StatusUnauthorized)
		return
	}

	err := handler.userRepository.DeleteUser(req.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
