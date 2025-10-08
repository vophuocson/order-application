package main

import (
	"fmt"
	"net/http"
	"time"
	"user-domain/config"
	database "user-domain/db"
	router "user-domain/infrastructure/http"
	"user-domain/infrastructure/logger"
)

type Server struct {
	httpServer *http.Server
}

func main() {
	cfg := config.LoadConfig()
	db, err := database.NewDatabase(cfg)
	if err != nil {
		panic(err.Error())
	}
	gorm, err := db.NewGorm()
	if err != nil {
		panic(err.Error())
	}
	logger := logger.NewLogger()
	r := router.BuildRouter(gorm, logger)
	s := Server{
		httpServer: &http.Server{
			Handler:      r,
			Addr:         fmt.Sprintf(":%s", cfg.ApiPort),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
	err = s.httpServer.ListenAndServe()
	if err != nil {
		panic("error running server")
	}
}
