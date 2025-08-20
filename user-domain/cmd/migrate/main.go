package main

import (
	"database/sql"
	"fmt"
	"log"
	"user-domain/repository/env"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func connectDatabase() (*sql.DB, error) {
	cfg := env.NewConfig()
	err := cfg.Load("../../../.databse_env")
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.GetHostName(), cfg.GetPort(), cfg.GetUser(), cfg.GetDBName(), cfg.GetPassword(), cfg.GetSSLMode())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, err := connectDatabase()
	if err != nil {
		panic("error connects database")
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./ddl",
		"postgres", driver)
	if err != nil {
		log.Fatal(err.Error())

	}
	m.Up()
}
