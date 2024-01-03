package ping

import (
	"encoding/json"
	pingService "github.com/anoriar/gophermart/internal/gophermart/shared/services/ping"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

type PingHandler struct {
	pingService pingService.PingServiceInterface
}

func NewPingHandler(pingService pingService.PingServiceInterface) *PingHandler {
	return &PingHandler{pingService: pingService}
}

func (handler *PingHandler) Ping(w http.ResponseWriter, req *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(req.Context(), "PingHandler::Ping")
	defer span.Finish()

	response := handler.pingService.Ping(ctx)
	responseBody, err := json.Marshal(response)

	span.LogFields(
		log.Object("response", response),
	)

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
