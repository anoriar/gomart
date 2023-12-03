package app

import (
	"github.com/anoriar/gophermart/internal/gophermart/balance/repository/balance"
	balanceServicePkg "github.com/anoriar/gophermart/internal/gophermart/balance/services/balance"
	"github.com/anoriar/gophermart/internal/gophermart/balance/services/withdraw"
	"github.com/anoriar/gophermart/internal/gophermart/order/services"
	"github.com/anoriar/gophermart/internal/gophermart/order/services/fetcher"
	"github.com/anoriar/gophermart/internal/gophermart/shared/app/db"
	"github.com/anoriar/gophermart/internal/gophermart/shared/config"
	"github.com/anoriar/gophermart/internal/gophermart/shared/services/ping"
	"github.com/anoriar/gophermart/internal/gophermart/shared/services/validator/idvalidator"
	"github.com/anoriar/gophermart/internal/gophermart/user/repository"
	"github.com/anoriar/gophermart/internal/gophermart/user/services/auth"
	"go.uber.org/zap"
)

type App struct {
	Config            *config.Config
	Logger            *zap.Logger
	Database          *db.Database
	PingService       ping.PingServiceInterface
	AuthService       auth.AuthServiceInterface
	OrderService      services.OrderServiceInterface
	OrderFetchService fetcher.OrderFetchServiceInterface
	UserRepository    repository.UserRepositoryInterface
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
	orderService services.OrderServiceInterface,
	orderFetchService fetcher.OrderFetchServiceInterface,
	userRepository repository.UserRepositoryInterface,
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
