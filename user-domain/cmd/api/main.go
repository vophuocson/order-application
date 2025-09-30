package main

import (
	"net/http"
	"user-domain/controler/handler"
	usercontroler "user-domain/controler/user"
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
	handlerStrict := handler.NewStrictHandler(userControler, nil)
	handler := handler.Handler(handlerStrict)
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
