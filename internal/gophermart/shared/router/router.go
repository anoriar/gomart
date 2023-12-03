package router

import (
	"github.com/anoriar/gophermart/internal/gophermart/balance/handlers/balance"
	"github.com/anoriar/gophermart/internal/gophermart/balance/handlers/getwithdrawals"
	"github.com/anoriar/gophermart/internal/gophermart/balance/handlers/withdraw"
	"github.com/anoriar/gophermart/internal/gophermart/order/handlers/getorders"
	"github.com/anoriar/gophermart/internal/gophermart/order/handlers/loadorder"
	"github.com/anoriar/gophermart/internal/gophermart/shared/app"
	"github.com/anoriar/gophermart/internal/gophermart/shared/handlers/ping"
	"github.com/anoriar/gophermart/internal/gophermart/shared/middleware/auth"
	"github.com/anoriar/gophermart/internal/gophermart/shared/middleware/compress"
	"github.com/anoriar/gophermart/internal/gophermart/shared/middleware/logger"
	"github.com/anoriar/gophermart/internal/gophermart/user/handlers/delete"
	"github.com/anoriar/gophermart/internal/gophermart/user/handlers/login"
	"github.com/anoriar/gophermart/internal/gophermart/user/handlers/register"
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
	userDeleteHandler     *delete.DeleteHandler
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
		userDeleteHandler:     delete.NewDeleteHandler(app.UserRepository),
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
	router.With(r.authMiddleware.Auth).Delete("/api/user", r.userDeleteHandler.Delete)

	return router
}
