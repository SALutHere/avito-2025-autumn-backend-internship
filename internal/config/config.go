package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	cfg  *Config
	once sync.Once
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	Postgres   `yaml:"postgres"`
}

type HTTPServer struct {
	HTTPPort         int           `yaml:"port" env-default:"8080"`
	HTTPReadTimeout  time.Duration `yaml:"read_timeout" env-default:"4s"`
	HTTPWriteTimeout time.Duration `yaml:"write_timeout" env-default:"4s"`
	HTTPIdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Postgres struct {
	PGHost     string        `env:"POSTGRES_HOST" env-default:"db"`
	PGPort     int           `env:"POSTGRES_PORT" env-default:"5432"`
	PGUser     string        `env:"POSTGRES_USER" env-default:"postgres"`
	PGPassword string        `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	PGDBName   string        `env:"POSTGRES_DB_NAME" env-default:"pr_service"`
	PGTimeout  time.Duration `yaml:"timeout" env-default:"4s"`
}

func Load(configPath string) *Config {
	once.Do(func() {
		if configPath == "" {
			log.Fatal("config path is not set")
		}

		if _, err := os.Stat(configPath); err != nil {
			log.Fatalf("config path does not exist: %s", configPath)
		}

		cfg = &Config{}

		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			log.Fatalf("can not read config: %v", err)
		}
	})

	return cfg
}

func C() *Config {
	return cfg
}

func (cfg *Config) PostgresURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PGUser,
		cfg.PGPassword,
		cfg.PGHost,
		cfg.PGPort,
		cfg.PGDBName,
	)
}
