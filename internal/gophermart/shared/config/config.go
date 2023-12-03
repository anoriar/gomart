package config

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel             string `env:"LOG_LEVEL"`
	JwtSecretKey         string `env:"JWT_SECRET_KEY"`
}

func NewConfig() *Config {
	return &Config{
		LogLevel: "debug",
	}
}
