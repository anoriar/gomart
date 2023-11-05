package ping

import (
	"encoding/json"
	pingService "github.com/anoriar/gophermart/internal/gophermart/services/ping"
	"net/http"
)

type PingHandler struct {
	pingService pingService.PingServiceInterface
}

func NewPingHandler(pingService pingService.PingServiceInterface) *PingHandler {
	return &PingHandler{pingService: pingService}
}

func (handler *PingHandler) Ping(w http.ResponseWriter, req *http.Request) {
	response := handler.pingService.Ping()
	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Json marshal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBody)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
