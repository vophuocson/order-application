package main

import (
	"net/http"
	database "user-domain/db"
	userapplication "user-domain/internal/application/controller/user"
	applicationlogger "user-domain/internal/application/logger"
	applicationoutbound "user-domain/internal/application/outbound"
	userrepository "user-domain/internal/application/repository/user"
	userdomain "user-domain/internal/domain/user"
	"user-domain/pkg/handler"
	"user-domain/pkg/logger"
	userpersistence "user-domain/pkg/persistence/user"

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
	r.Route("/api/v1", func(r chi.Router) {
		buildUserSubRouter(r, db, logger)
	})
}
