package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `env:"TODO_ENV" envDefault:"dev"`
	Port        string `env:"PORT" envDefault:":8080"`
	DSN         string `env:"DSN" envDefault:"user:password@tcp(localhost:3306)/dbname"`
	ExternalAPI string `env:"EXTERNAL_API" envDefault:"https://geoapi.heartrails.com/api/json?method=searchByPostal&postal="`
}

func New() (*Config, error) {
	// .env ファイルを読み込む
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
