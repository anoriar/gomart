package main

import (
	"flag"
	"github.com/anoriar/gophermart/internal/gophermart/shared/config"
)

func parseFlags(config *config.Config) {
	flag.StringVar(&config.RunAddress, "a", "localhost:8080", "Run address")
	flag.StringVar(&config.DatabaseURI, "d", "", "Database dsn")
	flag.StringVar(&config.AccrualSystemAddress, "r", "http://localhost:8080", "Accrual system address")
	flag.StringVar(&config.JwtSecretKey, "s", "", "JWT secret key")

	flag.Parse()
}
