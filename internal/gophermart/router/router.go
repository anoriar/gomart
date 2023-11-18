package router

import (
	"github.com/anoriar/gophermart/internal/gophermart/app"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/balance"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/getorders"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/getwithdrawals"
	"github.com/anoriar/gophermart/internal/gophermart/handlers/loadorder"
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
	loggerMiddleware      *logger.LoggerMiddleware
	compressMiddleware    *compress.CompressMiddleware
	authMiddleware        *auth.AuthMiddleware
	pingHandler           *ping.PingHandler
	registerHandler       *register.RegisterHandler
	loginHandler          *login.LoginHandler
	loadOrderHandler      *loadorder.LoadOrderHandler
	getOrdersHandler      *getorders.GetOrdersHandler
	balanceHandler        *balance.BalanceHandler
	withdrawHandler       *withdraw.WithdrawHandler
	getWithdrawalsHandler *getwithdrawals.GetWithdrawalsHandler
}

func NewRouter(app *app.App) *Router {
	return &Router{
		loggerMiddleware:      logger.NewLoggerMiddleware(app.Logger),
		compressMiddleware:    compress.NewCompressMiddleware(),
		authMiddleware:        auth.NewAuthMiddleware(app.AuthService),
		pingHandler:           ping.NewPingHandler(app.PingService),
		registerHandler:       register.NewRegisterHandler(app.AuthService),
		loginHandler:          login.NewLoginHandler(app.AuthService),
		loadOrderHandler:      loadorder.NewLoadOrderHandler(app.OrderService),
		getOrdersHandler:      getorders.NewGetOrdersHandler(app.OrderService),
		balanceHandler:        balance.NewBalanceHandler(app.BalanceService),
		withdrawHandler:       withdraw.NewWithdrawHandler(app.WithdrawService),
		getWithdrawalsHandler: getwithdrawals.NewGetWithdrawalsHandler(app.WithdrawService),
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
	router.With(r.authMiddleware.Auth).Get("/api/user/withdrawals", r.getWithdrawalsHandler.GetWithdrawals)

	return router
}
