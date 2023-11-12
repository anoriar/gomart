package app

import (
	dbPkg "github.com/anoriar/gophermart/internal/gophermart/app/db"
	loggerPkg "github.com/anoriar/gophermart/internal/gophermart/app/logger"
	"github.com/anoriar/gophermart/internal/gophermart/config"
	"github.com/anoriar/gophermart/internal/gophermart/repository/accrual"
	orderPkg "github.com/anoriar/gophermart/internal/gophermart/repository/order"
	"github.com/anoriar/gophermart/internal/gophermart/repository/user"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth"
	"github.com/anoriar/gophermart/internal/gophermart/services/order"
	"github.com/anoriar/gophermart/internal/gophermart/services/ping"
	"net/http"
)

func InitializeApp(conf *config.Config) (*App, error) {
	db, err := dbPkg.InitializeDatabase(conf.DatabaseURI)
	if err != nil {
		return nil, err
	}
	logger, err := loggerPkg.Initialize(conf.LogLevel)
	if err != nil {
		return nil, err
	}
	userRepository := user.NewUserRepository(db)
	orderRepository := orderPkg.NewOrderRepository(db)
	accrualRepository := accrual.NewAccrualRepository(&http.Client{}, conf.AccrualSystemAddress)
	authService := auth.InitializeAuthService(conf, userRepository, logger)
	orderService := order.NewOrderService(orderRepository, accrualRepository, logger)

	return NewApp(conf, logger, db, ping.NewPingService(db), authService, orderService, userRepository), nil
}
