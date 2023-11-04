package app

import (
	"github.com/anoriar/gophermart/internal/gophermart/app/db"
	"github.com/anoriar/gophermart/internal/gophermart/config"
	"github.com/anoriar/gophermart/internal/gophermart/repository/user"
	"go.uber.org/zap"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	Database       *db.Database
	UserRepository user.UserRepositoryInterface
}

func NewApp(config *config.Config, logger *zap.Logger, database *db.Database, userRepository user.UserRepositoryInterface) *App {
	return &App{Config: config, Logger: logger, Database: database, UserRepository: userRepository}
}

func (app *App) Close() {
	app.Database.Close()
	app.Logger.Sync()
}
