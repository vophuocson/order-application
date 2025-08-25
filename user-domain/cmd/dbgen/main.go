package main

import (
	database "user-domain/pkg/db"

	"github.com/pkg/errors"
	"gorm.io/gen"
)

func main() {
	cfg := gen.Config{
		OutPath:          "../../repository/dao",
		FieldWithTypeTag: true,
		Mode:             gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	}
	g := gen.NewGenerator(cfg)

	db, err := database.NewDatabase()
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
