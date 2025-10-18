package main

import (
	"log"
	"user-domain/infrastructure/config"
	"user-domain/infrastructure/database"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbConfig := config.LoadConfig()
	conn, err := database.NewDatabaseConection(dbConfig)
	if err != nil {
		panic("error connects database")
	}
	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./cmd/migrate/ddl",
		"postgres", driver)
	if err != nil {
		log.Fatal(err.Error())

	}
	m.Up()
}
