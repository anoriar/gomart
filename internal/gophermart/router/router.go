package router

import (
	"github.com/anoriar/gophermart/internal/gophermart/app"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/ping"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	app         *app.App
	pingHandler *ping.PingHandler
}

func NewRouter(app *app.App) *Router {
	return &Router{
		app:         app,
		pingHandler: ping.NewPingHandler(app),
	}
}

func (r *Router) Route() chi.Router {
	router := chi.NewRouter()

	router.Get("/ping", r.pingHandler.Ping)

	return router
}
