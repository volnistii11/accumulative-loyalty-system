package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type ParserGetter interface {
	Parse() error
	GetStorageDSN() string
	GetHTTPServerAddress() string
	GetAccrualSystemAddress() string
}

func New() ParserGetter {
	return &Config{
		StorageDSN:           "",
		HTTPServerAddress:    "",
		AccrualSystemAddress: "",
	}
}

type Config struct {
	StorageDSN           string `yaml:"storage_dsn"`
	HTTPServerAddress    string `yaml:"http_server_address"`
	AccrualSystemAddress string `yaml:"accrual_system_address"`
}

func (cfg *Config) Parse() error {
	defaultValues, err := parseDefaultValues()
	if err != nil {
		return errors.Wrap(err, "Config.parse")
	}

	flag.StringVar(&cfg.StorageDSN, "d", defaultValues.StorageDSN, "repository URI")
	flag.StringVar(&cfg.HTTPServerAddress, "a", defaultValues.HTTPServerAddress, "http server address")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", defaultValues.AccrualSystemAddress, "accrual system address")
	flag.Parse()

	if envStorageDSN := os.Getenv("DATABASE_URI"); envStorageDSN != "" {
		cfg.StorageDSN = envStorageDSN
	}
	if envHTTPServerAddress := os.Getenv("RUN_ADDRESS"); envHTTPServerAddress != "" {
		cfg.HTTPServerAddress = envHTTPServerAddress
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		cfg.AccrualSystemAddress = envAccrualSystemAddress
	}

	if cfg.StorageDSN == "" {
		return errors.New("repository DSN is empty")
	}
	if cfg.HTTPServerAddress == "" {
		return errors.New("http server address is empty")
	}
	if cfg.AccrualSystemAddress == "" {
		return errors.New("accrual system address is empty")
	}
	return nil
}

func (cfg *Config) GetStorageDSN() string {
	return cfg.StorageDSN
}

func (cfg *Config) GetHTTPServerAddress() string {
	return cfg.HTTPServerAddress
}

func (cfg *Config) GetAccrualSystemAddress() string {
	return cfg.AccrualSystemAddress
}

func parseDefaultValues() (*Config, error) {
	configPath := "./config/local.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New("config file with default values does not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, errors.New("cannot read config file with default values")
	}

	return &cfg, nil
}
