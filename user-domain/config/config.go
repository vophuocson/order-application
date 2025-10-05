package config

import (
	"os"
)

type Config struct {
	PostgresDatabase string
	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresSSLMode  string
	PostgresPort     string
	ApiPort          string
}

func LoadConfig() *Config {
	cfg := &Config{
		PostgresDatabase: os.Getenv("SECRET_POSTGRES_DATABASE"),
		PostgresHost:     os.Getenv("SECRET_POSTGRES_HOSTNAME"),
		PostgresUser:     os.Getenv("SECRET_POSTGRES_USER"),
		PostgresPassword: os.Getenv("SECRET_POSTGRES_PASSWORD"),
		PostgresSSLMode:  os.Getenv("SECRET_POSTGRES_SSL_MODE"),
		PostgresPort:     os.Getenv("SECRET_POSTGRES_PORT"),
		ApiPort:          os.Getenv("API_PORT"),
	}

	return cfg
}
