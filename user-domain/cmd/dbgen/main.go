package main

import (
	database "user-domain/db"
	"user-domain/infrastructure/config"

	"github.com/pkg/errors"
	"gorm.io/gen"
)

func main() {
	cfg := gen.Config{
		OutPath:          "./infrastructure/persistence/postgres/dao",
		FieldWithTypeTag: true,
		Mode:             gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	}
	g := gen.NewGenerator(cfg)

	dbConfig := config.LoadConfig()
	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		panic(errors.Wrap(err, "connect database failed"))
	}
	gorm, err := db.NewGorm()
	if err != nil {
		panic(errors.Wrap(err, "init gorm failed"))
	}
	g.WithDataTypeMap(map[string]func(detailType string) (dataType string){
		"uuid": func(detailType string) (dataType string) {
			return "string"
		},
	})

	g.UseDB(gorm)
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}
