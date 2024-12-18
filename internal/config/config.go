package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Server struct {
		Host string `env:"SERVER_HOST" envDefault:"localhost"`
		Port string `env:"SERVER_PORT" envDefault:"7540"`
	}

	DB struct {
		Path string `env:"TODO_DBFILE" envDefault:""`
		Name string `env:"DB_NAME" envDefault:"scheduler"`
	}

	User struct {
		Password string `env:"TODO_PASSWORD" envDefault:""`
	}

	JWT struct {
		TTL    int    `env:"JWT_TTL" envDefault:"8"`
		Secret string `env:"JWT_SECRET" envDefault:"secret"`
	}
)

type Config struct {
	Server
	DB
	User
	JWT
}

// LoadConfig загружает конфигурацию из.env файла и возвращает её экземпляр.
func LoadConfig() (*Config, error) {
	// Загружаем конфигурацию из .env файла
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
