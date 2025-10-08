package main

import (
	"fmt"
	"net/http"
	"time"
	"user-domain/config"
	database "user-domain/db"
	"user-domain/infrastructure/logger"
	userpersistence "user-domain/infrastructure/persistence/user"
	userapplication "user-domain/internal/application/controller/user"
	applicationlogger "user-domain/internal/application/logger"
	applicationoutbound "user-domain/internal/application/outbound"
	userrepository "user-domain/internal/application/repository/user"
	userdomain "user-domain/internal/domain/user"
	"user-domain/pkg/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
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
	r := chi.NewRouter()
	buidRouter(r, gorm, logger)
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

func buildUserSubRouter(r chi.Router, db *gorm.DB, loggerOutbound applicationoutbound.Logger) {
	userPersistence := userpersistence.NewUserRepo(db)
	userRepo := userrepository.NewUserRepo(userPersistence)

	loggerOutport := applicationlogger.NewLogger(loggerOutbound)
	userService := userdomain.NewUserService(userRepo, loggerOutport)

	userControler := userapplication.NewUserControler(userService, loggerOutbound)
	handler.HandlerWithOptions(userControler, handler.ChiServerOptions{
		BaseRouter: r,
	})
}

func buidRouter(r chi.Router, db *gorm.DB, logger applicationoutbound.Logger) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Route("/api/v1", func(r chi.Router) {
		buildUserSubRouter(r, db, logger)
	})
}
