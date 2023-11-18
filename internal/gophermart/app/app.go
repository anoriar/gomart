package app

import (
	"github.com/anoriar/gophermart/internal/gophermart/app/db"
	"github.com/anoriar/gophermart/internal/gophermart/config"
	"github.com/anoriar/gophermart/internal/gophermart/repository/balance"
	"github.com/anoriar/gophermart/internal/gophermart/repository/user"
	"github.com/anoriar/gophermart/internal/gophermart/services/auth"
	balanceServicePkg "github.com/anoriar/gophermart/internal/gophermart/services/balance"
	"github.com/anoriar/gophermart/internal/gophermart/services/order"
	"github.com/anoriar/gophermart/internal/gophermart/services/order/fetcher"
	"github.com/anoriar/gophermart/internal/gophermart/services/ping"
	"github.com/anoriar/gophermart/internal/gophermart/services/validator/idvalidator"
	"github.com/anoriar/gophermart/internal/gophermart/services/withdraw"
	"go.uber.org/zap"
)

type App struct {
	Config            *config.Config
	Logger            *zap.Logger
	Database          *db.Database
	PingService       ping.PingServiceInterface
	AuthService       auth.AuthServiceInterface
	OrderService      order.OrderServiceInterface
	OrderFetchService fetcher.OrderFetchServiceInterface
	UserRepository    user.UserRepositoryInterface
	BalanceRepository balance.BalanceRepositoryInterface
	BalanceService    balanceServicePkg.BalanceServiceInterface
	IDValidator       idvalidator.IDValidatorInterface
	WithdrawService   withdraw.WithdrawServiceInterface
}

func NewApp(
	config *config.Config,
	logger *zap.Logger,
	database *db.Database,
	pingService ping.PingServiceInterface,
	authService auth.AuthServiceInterface,
	orderService order.OrderServiceInterface,
	orderFetchService fetcher.OrderFetchServiceInterface,
	userRepository user.UserRepositoryInterface,
	balanceRepository balance.BalanceRepositoryInterface,
	balanceService balanceServicePkg.BalanceServiceInterface,
	idValidator idvalidator.IDValidatorInterface,
	withdrawService withdraw.WithdrawServiceInterface,

) *App {
	return &App{
		Config:            config,
		Logger:            logger,
		Database:          database,
		PingService:       pingService,
		AuthService:       authService,
		OrderService:      orderService,
		OrderFetchService: orderFetchService,
		UserRepository:    userRepository,
		BalanceRepository: balanceRepository,
		BalanceService:    balanceService,
		IDValidator:       idValidator,
		WithdrawService:   withdrawService,
	}
}

func (app *App) Close() {
	app.Database.Close()
	app.Logger.Sync()
}
