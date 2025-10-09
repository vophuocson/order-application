package main

import (
	"log"
	database "user-domain/db"
	"user-domain/infrastructure/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbConfig := config.LoadConfig()
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		panic("error connects database")
	}
	driver, err := postgres.WithInstance(db.GetConnect(), &postgres.Config{})
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
