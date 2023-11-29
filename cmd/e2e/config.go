package e2e

type Config struct {
	GophermartBin     string `env:"GOPHERMART_BIN"`
	GophermartAddress string `env:"GOPHERMART_ADDRESS"`

	DatabaseURI string `env:"DATABASE_URI"`

	AccrualSystemBin     string `env:"ACCRUAL_SYSTEM_BIN"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() *Config {
	return &Config{}
}
