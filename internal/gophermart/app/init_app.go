package app

import (
	"context"
	dbPkg "github.com/anoriar/gophermart/internal/gophermart/app/db"
	loggerPkg "github.com/anoriar/gophermart/internal/gophermart/app/logger"
	"github.com/anoriar/gophermart/internal/gophermart/config"
	"github.com/anoriar/gophermart/internal/gophermart/processors/bus"
	orderProcessorPkg "github.com/anoriar/gophermart/internal/gophermart/processors/order"
	"github.com/anoriar/gophermart/internal/gophermart/repository/accrual"
	orderPkg "github.com/anoriar/gophermart/internal/gophermart/repository/order"
	"github.com/anoriar/gophermart/internal/gophermart/repository/user"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth"
	"github.com/anoriar/gophermart/internal/gophermart/services/order"
	"github.com/anoriar/gophermart/internal/gophermart/services/order/fetcher"
	"github.com/anoriar/gophermart/internal/gophermart/services/ping"
	"net/http"
)

func InitializeApp(ctx context.Context, conf *config.Config) (*App, error) {
	db, err := dbPkg.InitializeDatabase(conf.DatabaseURI)
	if err != nil {
		return nil, err
	}
	logger, err := loggerPkg.Initialize(conf.LogLevel)
	if err != nil {
		return nil, err
	}

	messageBus := bus.NewMessageBus()

	userRepository := user.NewUserRepository(db)
	orderRepository := orderPkg.NewOrderRepository(db)
	accrualRepository := accrual.NewAccrualRepository(&http.Client{}, conf.AccrualSystemAddress)
	authService := auth.InitializeAuthService(conf, userRepository, logger)
	orderFetchService := fetcher.NewOrderFetchService(accrualRepository)
	orderService := order.NewOrderService(orderRepository, orderFetchService, messageBus, logger)

	//запуск горутин
	orderProcessorPkg.NewOrderProcessor(ctx, orderService, logger, messageBus.OrderProcessChan)

	return NewApp(conf, logger, db, ping.NewPingService(db), authService, orderService, orderFetchService, userRepository), nil
}
