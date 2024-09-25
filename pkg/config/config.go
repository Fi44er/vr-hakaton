package config

import (
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const (
	ProductionEnv   = "development"
	DatabaseTimeout = 5 * time.Second
)

type Schema struct {
	Enviroment  string `env:"enviroment"`
	HttpPort    int    `env:"http_port"`
	AuthSecret  string `env:"auth_secret"`
	DatabaseURI string `env:"database_uri"`
}

var cfg Schema

func LoadConfig() *Schema {
	_, filename, _, _ := runtime.Caller(0)
	cirrentDir := filepath.Dir(filename)

	if err := godotenv.Load(filepath.Join(cirrentDir, "config.yaml")); err != nil {
		log.Printf("Error on load configuration file, error: %v", err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error on parsing configuration file, error: %v", err)
	}

	return &cfg
}

func GetConfig() *Schema {
	return &cfg
}
