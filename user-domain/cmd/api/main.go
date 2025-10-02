package main

import (
	"net/http"
	"user-domain/controller/handler"
	usercontroler "user-domain/controller/user"
	"user-domain/internal/service"
	database "user-domain/pkg/db"
	"user-domain/pkg/logger"
	userrepo "user-domain/repository/user"
)

type Server struct {
	httpServer *http.Server
}

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		panic(err.Error())
	}
	gorm, err := db.NewGorm()
	if err != nil {
		panic(err.Error())
	}
	userRepo := userrepo.NewUserRepo(gorm)
	logger := logger.NewLogger()
	userService := service.NewUserService(userRepo, logger)
	userControler := usercontroler.NewUserControler(userService, logger)
	handler := handler.Handler(userControler)
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
