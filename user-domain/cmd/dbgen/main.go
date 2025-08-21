package main

import (
	"fmt"
	"user-domain/env"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	cfg := gen.Config{
		OutPath:          "../../repository/dao",
		FieldWithTypeTag: true,
		Mode:             gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	}
	g := gen.NewGenerator(cfg)
	db, err := connectDatabase()
	if err != nil {
		panic(errors.Wrap(err, "Connect database failed"))
	}
	g.UseDB(db)
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}

func connectDatabase() (*gorm.DB, error) {
	cfg := env.NewConfig()
	err := cfg.Load("../../../.databse_env")
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.GetHostName(), cfg.GetPort(), cfg.GetUser(), cfg.GetDBName(), cfg.GetPassword(), cfg.GetSSLMode())
	gormConfig := &gorm.Config{}
	return gorm.Open(postgres.Open(dsn), gormConfig)
}
