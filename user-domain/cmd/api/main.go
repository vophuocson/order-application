package main

import (
	"fmt"
	"net/http"
	"time"
	"user-domain/infrastructure/config"
	"user-domain/infrastructure/database"
	router "user-domain/infrastructure/http"
	"user-domain/infrastructure/logger"
)

type Server struct {
	httpServer *http.Server
}

func main() {
	cfg := config.LoadConfig()
	logger := logger.NewLogger()
	gorm, err := database.NewGorm(cfg, logger)
	if err != nil {
		panic(err.Error())
	}
	// flush buffer before exiting
	defer logger.Sync()
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
