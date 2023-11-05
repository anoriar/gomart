package router

import (
	"github.com/anoriar/gophermart/internal/gophermart/app"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/ping"
	"github.com/anoriar/gophermart/internal/gophermart/middleware/compress"
	"github.com/anoriar/gophermart/internal/gophermart/middleware/logger"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	app                *app.App
	loggerMiddleware   *logger.LoggerMiddleware
	compressMiddleware *compress.CompressMiddleware
	pingHandler        *ping.PingHandler
}

func NewRouter(app *app.App) *Router {
	return &Router{
		app:                app,
		loggerMiddleware:   logger.NewLoggerMiddleware(app.Logger),
		compressMiddleware: compress.NewCompressMiddleware(),
		pingHandler:        ping.NewPingHandler(app.PingService),
	}
}

func (r *Router) Route() chi.Router {
	router := chi.NewRouter()

	router.Use(r.loggerMiddleware.Log)
	router.Use(r.compressMiddleware.Compress)

	router.Get("/ping", r.pingHandler.Ping)

	return router
}
