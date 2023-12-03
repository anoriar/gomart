package app

import (
	"context"
	balanceRepositoryPkg "github.com/anoriar/gophermart/internal/gophermart/balance/repository/balance"
	withdrawalRepositoryPkg "github.com/anoriar/gophermart/internal/gophermart/balance/repository/withdrawal"
	"github.com/anoriar/gophermart/internal/gophermart/balance/services/balance"
	"github.com/anoriar/gophermart/internal/gophermart/balance/services/withdraw"
	orderProcessorPkg "github.com/anoriar/gophermart/internal/gophermart/order/processors/orderprocess"
	orderprocess "github.com/anoriar/gophermart/internal/gophermart/order/processors/syncfailed"
	orderPkg "github.com/anoriar/gophermart/internal/gophermart/order/repository"
	"github.com/anoriar/gophermart/internal/gophermart/order/repository/accrual"
	"github.com/anoriar/gophermart/internal/gophermart/order/services"
	"github.com/anoriar/gophermart/internal/gophermart/order/services/fetcher"
	dbPkg "github.com/anoriar/gophermart/internal/gophermart/shared/app/db"
	loggerPkg "github.com/anoriar/gophermart/internal/gophermart/shared/app/logger"
	"github.com/anoriar/gophermart/internal/gophermart/shared/config"
	"github.com/anoriar/gophermart/internal/gophermart/shared/services/bus"
	"github.com/anoriar/gophermart/internal/gophermart/shared/services/ping"
	"github.com/anoriar/gophermart/internal/gophermart/shared/services/validator/idvalidator"
	"github.com/anoriar/gophermart/internal/gophermart/user/repository"
	"github.com/anoriar/gophermart/internal/gophermart/user/services/auth"
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

	messageBus := bus.NewMessageBus(logger)

	userRepository := repository.NewUserRepository(db)
	orderRepository := orderPkg.NewOrderRepository(db)
	accrualRepository := accrual.NewAccrualRepository(&http.Client{}, conf.AccrualSystemAddress)
	withdrawalRepository := withdrawalRepositoryPkg.NewWithdrawalRepository(db)
	balanceRepository := balanceRepositoryPkg.NewBalanceRepository(db, withdrawalRepository)

	idValidator := idvalidator.NewLuhnValidator()

	authService := auth.InitializeAuthService(conf, userRepository, logger)
	orderFetchService := fetcher.NewOrderFetchService(accrualRepository)
	balanceService := balance.NewBalanceService(balanceRepository, logger)
	orderService := services.NewOrderService(orderRepository, orderFetchService, balanceService, messageBus, idValidator, logger)

	withdrawService := withdraw.NewWithdrawService(withdrawalRepository, balanceService, idValidator, logger)

	//запуск горутин
	orderProcessorPkg.NewOrderProcessor(ctx, orderService, logger, messageBus.OrderProcessChan)
	orderprocess.NewOrderSyncFailedProcessor(ctx, orderService, logger, messageBus)

	return NewApp(
		conf,
		logger,
		db,
		ping.NewPingService(db),
		authService,
		orderService,
		orderFetchService,
		userRepository,
		balanceRepository,
		balanceService,
		idValidator,
		withdrawService,
	), nil
}
