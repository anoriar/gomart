package router

import (
	"github.com/anoriar/gophermart/internal/gophermart/app"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/ping"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/register"
	"github.com/anoriar/gophermart/internal/gophermart/middleware/compress"
	"github.com/anoriar/gophermart/internal/gophermart/middleware/logger"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	loggerMiddleware   *logger.LoggerMiddleware
	compressMiddleware *compress.CompressMiddleware
	pingHandler        *ping.PingHandler
	registerHandler    *register.RegisterHandler
}

func NewRouter(app *app.App) *Router {
	return &Router{
		loggerMiddleware:   logger.NewLoggerMiddleware(app.Logger),
		compressMiddleware: compress.NewCompressMiddleware(),
		pingHandler:        ping.NewPingHandler(app.PingService),
		registerHandler:    register.NewRegisterHandler(app.AuthService),
	}
}

func (r *Router) Route() chi.Router {
	router := chi.NewRouter()

	router.Use(r.loggerMiddleware.Log)
	router.Use(r.compressMiddleware.Compress)

	router.Get("/ping", r.pingHandler.Ping)
	router.Post("/api/user/register", r.registerHandler.Register)

	return router
}
