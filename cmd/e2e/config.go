package e2e

type Config struct {
	GophermartBin        string `env:"GOPHERMART_BIN"`
	GophermartRunAddress string `env:"GOPHERMART_RUN_ADDRESS"`
	GophermartAddress    string `env:"GOPHERMART_ADDRESS"`

	DatabaseURI string `env:"DATABASE_URI"`

	AccrualSystemBin        string `env:"ACCRUAL_SYSTEM_BIN"`
	AccrualSystemRunAddress string `env:"ACCRUAL_SYSTEM_RUN_ADDRESS"`
	AccrualSystemAddress    string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() *Config {
	return &Config{}
}
