package main

import (
	"context"
	appPkg "github.com/anoriar/gophermart/internal/gophermart/shared/app"
	"github.com/anoriar/gophermart/internal/gophermart/shared/config"
	"github.com/anoriar/gophermart/internal/gophermart/shared/router"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	conf := config.NewConfig()
	parseFlags(conf)

	err := env.Parse(conf)
	if err != nil {
		panic(err)
	}

	app, err := appPkg.InitializeApp(context.Background(), conf)
	if err != nil {
		panic(err)
	}
	defer app.Close()

	r := router.NewRouter(app)

	err = http.ListenAndServe(conf.RunAddress, r.Route())
	if err != nil {
		app.Logger.Fatal("Server exception", zap.String("exception", err.Error()))
		panic(err)
	}
}
