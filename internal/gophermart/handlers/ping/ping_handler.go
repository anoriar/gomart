package ping

import (
	"github.com/anoriar/gophermart/internal/gophermart/app"
	"net/http"
)

type PingHandler struct {
	app *app.App
}

func NewPingHandler(app *app.App) *PingHandler {
	return &PingHandler{app: app}
}

func (handler *PingHandler) Ping(w http.ResponseWriter, req *http.Request) {
	//TODO: ping

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
