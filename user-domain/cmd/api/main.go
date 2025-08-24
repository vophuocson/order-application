package main

import (
	"fmt"
	"net/http"
	"user-domain/api/handler"
	"user-domain/api/inbound"
	"user-domain/env"
	"user-domain/internal/service"
	userrepo "user-domain/repository/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	httpServer *http.Server
}

func main() {
	db, err := connectDatabase()
	if err != nil {
		panic("error connects database")
	}
	userRepo := userrepo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := inbound.NewUserHandler(userService)
	control := handler.NewStrictHandler(userHandler, nil)
	handler := handler.Handler(control)
	s := Server{
		httpServer: &http.Server{
			Handler: handler,
		},
	}
	err = s.httpServer.ListenAndServe()
	if err != nil {
		panic("error running server")
	}
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
