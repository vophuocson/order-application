package main

import (
	"user-domain/config"
	database "user-domain/db"

	"github.com/pkg/errors"
	"gorm.io/gen"
)

func main() {
	cfg := gen.Config{
		OutPath:          "./pgk/persistence/dao",
		FieldWithTypeTag: true,
		Mode:             gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	}
	g := gen.NewGenerator(cfg)

	dbConfig := config.LoadConfig()
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		panic(errors.Wrap(err, "Connect database failed"))
	}
	gorm, err := db.NewGorm()
	if err != nil {
		panic(errors.Wrap(err, "Connect database failed"))
	}

	g.UseDB(gorm)
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}
