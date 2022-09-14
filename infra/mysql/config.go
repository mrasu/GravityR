package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/mrasu/GravityR/lib"
	"os"
)

type Config struct {
	config    *mysql.Config
	paramText string
}

func NewConfigFromEnv() (*Config, error) {
	cfg := mysql.NewConfig()
	username, err := lib.GetEnv("DB_USERNAME")
	if err != nil {
		return nil, err
	}
	cfg.User = username

	cfg.Passwd = os.Getenv("DB_PASSWORD")
	protocol := os.Getenv("DB_PROTOCOL")
	if len(protocol) == 0 {
		protocol = "tcp"
	}
	cfg.Net = protocol

	address := os.Getenv("DB_ADDRESS")
	if len(address) > 0 {
		cfg.Addr = address
	}

	database, err := lib.GetEnv("DB_DATABASE")
	if err != nil {
		return nil, err
	}
	cfg.DBName = database

	paramText := os.Getenv("DB_PARAM_TEXT")

	return &Config{
		config:    cfg,
		paramText: paramText,
	}, nil
}

func (cfg *Config) ToDSN() string {
	dsn := cfg.config.FormatDSN()
	if len(cfg.paramText) > 0 {
		dsn += "?" + cfg.paramText
	}

	return dsn
}

func (cfg *Config) GetDBName() string {
	return cfg.config.DBName
}
