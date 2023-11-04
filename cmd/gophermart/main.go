package main

import (
	appPkg "github.com/anoriar/gophermart/internal/gophermart/app"
	"github.com/anoriar/gophermart/internal/gophermart/config"
	"github.com/anoriar/gophermart/internal/gophermart/router"
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

	app, err := appPkg.InitializeApp(conf)
	if err != nil {
		panic(err)
	}
	defer app.Close()

	r := router.NewRouter(app)

	if err != nil {
		app.Logger.Fatal("init app error", zap.String("error", err.Error()))
		panic(err)
	}

	err = http.ListenAndServe(conf.RunAddress, r.Route())
	if err != nil {
		app.Logger.Fatal("Server exception", zap.String("exception", err.Error()))
		panic(err)
	}
}
