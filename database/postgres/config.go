package postgres

import (
	"github.com/mrasu/GravityR/lib"
	"os"
	"strings"
)

type Config struct {
	user       string
	password   string
	host       string
	port       string
	DBName     string
	searchPath string
	paramText  string
}

func NewConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	username, err := lib.GetEnv("DB_USERNAME")
	if err != nil {
		return nil, err
	}
	cfg.user = username

	cfg.password = os.Getenv("DB_PASSWORD")

	cfg.host = os.Getenv("DB_HOST")
	cfg.port = os.Getenv("DB_PORT")

	database, err := lib.GetEnv("DB_DATABASE")
	if err != nil {
		return nil, err
	}
	cfg.DBName = database

	cfg.searchPath = os.Getenv("DB_SEARCH_PATH")

	cfg.paramText = os.Getenv("DB_PARAM_TEXT")

	return cfg, nil
}

// GetKVConnectionString returns libpq's connection string
// c.f) https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters
func (cfg *Config) GetKVConnectionString() string {
	var texts []string
	if cfg.user != "" {
		texts = append(texts, "user='"+cfg.user+"'")
	}
	if cfg.password != "" {
		texts = append(texts, "password='"+cfg.password+"'")
	}
	if cfg.host != "" {
		texts = append(texts, "host='"+cfg.host+"'")
	}
	if cfg.port != "" {
		texts = append(texts, "port='"+cfg.port+"'")
	}
	if cfg.DBName != "" {
		texts = append(texts, "dbname='"+cfg.DBName+"'")
	}
	if cfg.paramText != "" {
		texts = append(texts, "paramText='"+cfg.paramText+"'")
	}
	texts = append(texts, "sslmode=disable")

	return strings.Join(texts, " ")
}

func (cfg *Config) GetSearchPathOrPublic() string {
	if cfg.searchPath == "" {
		return "public"
	}

	return cfg.searchPath
}
