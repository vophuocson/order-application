package config

import (
	"os"

	"github.com/joho/godotenv"
)

type dbConfig struct {
	PostgresDatabase string `env:"SECRET_POSTGRES_DATABASE"`
	PostgresHostname string `env:"SECRET_POSTGRES_HOSTNAME"`
	PostgresPort     string `env:"SECRET_POSTGRES_PORT"`
	PostgresPassword string `env:"SECRET_POSTGRES_PASSWORD"`
	PostgresUser     string `env:"SECRET_POSTGRES_USER"`
	PostgresSSLMode  string `env:"SECRET_POSTGRES_SSL_MODE"`
}

func (db *dbConfig) GetDBName() string {
	return db.PostgresDatabase
}

func (db *dbConfig) GetHostName() string {
	return db.PostgresHostname
}

func (db *dbConfig) GetPort() string {
	return db.PostgresPort
}

func (db *dbConfig) GetPassword() string {
	return db.PostgresPassword
}

func (db *dbConfig) GetUser() string {
	return db.PostgresUser
}

func (db *dbConfig) GetSSLMode() string {
	return db.PostgresSSLMode
}

func (db *dbConfig) Load() error {
	i, err := os.Open(".env")
	if err != nil {
		return err
	}
	data, err := godotenv.Parse(i)
	if err != nil {
		return err
	}
	db.PostgresDatabase = data["SECRET_POSTGRES_DATABASE"]
	db.PostgresHostname = data["SECRET_POSTGRES_HOSTNAME"]
	db.PostgresPort = data["SECRET_POSTGRES_PORT"]
	db.PostgresPassword = data["SECRET_POSTGRES_PASSWORD"]
	db.PostgresUser = data["SECRET_POSTGRES_USER"]
	db.PostgresSSLMode = data["SECRET_POSTGRES_SSL_MODE"]
	return nil
}

type DBConfig interface {
	Load() error
	GetDBName() string
	GetHostName() string
	GetPort() string
	GetPassword() string
	GetUser() string
	GetSSLMode() string
}

func NewConfig() DBConfig {
	return &dbConfig{}
}
