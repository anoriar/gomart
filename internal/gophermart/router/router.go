package router

import (
	"github.com/anoriar/gophermart/internal/gophermart/app"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/balance"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/get_orders"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/load_order"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/login"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/ping"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/register"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/withdraw"
	"github.com/anoriar/gophermart/internal/gophermart/middleware/auth"
	"github.com/anoriar/gophermart/internal/gophermart/middleware/compress"
	"github.com/anoriar/gophermart/internal/gophermart/middleware/logger"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	loggerMiddleware   *logger.LoggerMiddleware
	compressMiddleware *compress.CompressMiddleware
	authMiddleware     *auth.AuthMiddleware
	pingHandler        *ping.PingHandler
	registerHandler    *register.RegisterHandler
	loginHandler       *login.LoginHandler
	loadOrderHandler   *load_order.LoadOrderHandler
	getOrdersHandler   *get_orders.GetOrdersHandler
	balanceHandler     *balance.BalanceHandler
	withdrawHandler    *withdraw.WithdrawHandler
}

func NewRouter(app *app.App) *Router {
	return &Router{
		loggerMiddleware:   logger.NewLoggerMiddleware(app.Logger),
		compressMiddleware: compress.NewCompressMiddleware(),
		authMiddleware:     auth.NewAuthMiddleware(app.AuthService),
		pingHandler:        ping.NewPingHandler(app.PingService),
		registerHandler:    register.NewRegisterHandler(app.AuthService),
		loginHandler:       login.NewLoginHandler(app.AuthService),
		loadOrderHandler:   load_order.NewLoadOrderHandler(app.OrderService),
		getOrdersHandler:   get_orders.NewGetOrdersHandler(app.OrderService),
		balanceHandler:     balance.NewBalanceHandler(app.BalanceService),
		withdrawHandler:    withdraw.NewWithdrawHandler(app.WithdrawService),
	}
}

func (r *Router) Route() chi.Router {
	router := chi.NewRouter()

	router.Use(r.loggerMiddleware.Log)
	router.Use(r.compressMiddleware.Compress)

	router.Get("/ping", r.pingHandler.Ping)
	router.Post("/api/user/register", r.registerHandler.Register)
	router.Post("/api/user/login", r.loginHandler.Login)
	router.With(r.authMiddleware.Auth).Post("/api/user/orders", r.loadOrderHandler.LoadOrder)
	router.With(r.authMiddleware.Auth).Get("/api/user/orders", r.getOrdersHandler.GetOrders)
	router.With(r.authMiddleware.Auth).Get("/api/user/balance", r.balanceHandler.GetUserBalance)
	router.With(r.authMiddleware.Auth).Post("/api/user/balance/withdraw", r.withdrawHandler.Withdraw)

	return router
}
