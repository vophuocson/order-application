package http

import (
	"user-domain/infrastructure/http/middleware"
	userpersistence "user-domain/infrastructure/persistence/user"
	usercontroller "user-domain/internal/application/controller/user"
	applicationlogger "user-domain/internal/application/logger"
	applicationoutbound "user-domain/internal/application/outbound"
	userrepository "user-domain/internal/application/repository/user"
	userdomain "user-domain/internal/domain/user"
	"user-domain/pkg/handler"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func BuildRouter(db *gorm.DB, logger applicationoutbound.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Route("/api/v1", func(r chi.Router) {
		buildUserSubRouter(r, db, logger)
	})
	return r
}

func buildUserSubRouter(r chi.Router, db *gorm.DB, loggerOutbound applicationoutbound.Logger) {
	userPersistence := userpersistence.NewUserRepo(db)
	userRepo := userrepository.NewUserRepo(userPersistence)

	loggerOutport := applicationlogger.NewLogger(loggerOutbound)
	userService := userdomain.NewUserService(userRepo, loggerOutport)

	userControler := usercontroller.NewUserControler(userService, loggerOutbound)
	handler.HandlerWithOptions(userControler, handler.ChiServerOptions{
		BaseRouter: r,
	})
}

// r := router.BuildRouter(gorm, logger)
// 	router "user-domain/infrastructure/http"
