package main

import (
	"net/http"
	"user-domain/controller/handler"
	usercontroler "user-domain/controller/user"
	"user-domain/internal/outport"
	"user-domain/internal/service"
	database "user-domain/pkg/db"
	"user-domain/pkg/logger"
	userrepo "user-domain/repository/user"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
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
	logger := logger.NewLogger()
	r := chi.NewRouter()
	buidRouter(r, gorm, logger)
	s := Server{
		httpServer: &http.Server{
			Handler: r,
		},
	}
	err = s.httpServer.ListenAndServe()
	if err != nil {
		panic("error running server")
	}
}

func buildUserSubRouter(r chi.Router, db *gorm.DB, logger outport.Logger) {
	userRepo := userrepo.NewUserRepo(db)
	userService := service.NewUserService(userRepo, logger)
	userControler := usercontroler.NewUserControler(userService, logger)
	handler.HandlerWithOptions(userControler, handler.ChiServerOptions{
		BaseRouter: r,
	})
}

func buidRouter(r chi.Router, db *gorm.DB, logger outport.Logger) {
	r.Route("/api/v1", func(r chi.Router) {
		buildUserSubRouter(r, db, logger)
	})
}
