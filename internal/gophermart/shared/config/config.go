package config

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	LogLevel             string `env:"LOG_LEVEL"`
	JwtSecretKey         string `env:"JWT_SECRET_KEY"`
	TracerServiceName    string `env:"TRACER_SERVICE_NAME"`
	TracerHeader         string `env:"TRACER_HEADER"`
}

func NewConfig() *Config {
	return &Config{
		LogLevel:          "debug",
		TracerServiceName: "gophermart",
		TracerHeader:      "x-traceid",
	}
}
