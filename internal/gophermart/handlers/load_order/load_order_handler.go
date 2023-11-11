package load_order

import (
	"fmt"
	"github.com/anoriar/gophermart/internal/gophermart/context"
	"net/http"
)

type LoadOrderHandler struct {
}

func NewLoadOrderHandler() *LoadOrderHandler {
	return &LoadOrderHandler{}
}

func (handler *LoadOrderHandler) LoadOrder(w http.ResponseWriter, req *http.Request) {
	userID := ""
	userIDCtxParam := req.Context().Value(context.UserIDContextKey)
	if userIDCtxParam != nil {
		userID = userIDCtxParam.(string)
	}
	fmt.Println(userID)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
