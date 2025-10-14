package http

import (
	"net/http"
	"user-domain/infrastructure/http/handler"
	"user-domain/infrastructure/http/middleware"
	userpersistence "user-domain/infrastructure/persistence/postgres/user"
	applicationparameter "user-domain/internal/application/controller/parameter"
	usercontroller "user-domain/internal/application/controller/user"
	applicationinbound "user-domain/internal/application/inbound"
	applicationlogger "user-domain/internal/application/logger"
	applicationoutbound "user-domain/internal/application/outbound"
	userrepository "user-domain/internal/application/repository/user"
	userdomain "user-domain/internal/domain/user"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func BuildRouter(db *gorm.DB, logger applicationoutbound.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.LoggingMiddleware(logger))
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

	handler.HandlerWithOptions(&userControllerWrap{UserApi: userControler}, handler.ChiServerOptions{
		BaseRouter: r,
	})
}

type userControllerWrap struct {
	applicationinbound.UserApi
}

func (cW *userControllerWrap) GetUsers(w http.ResponseWriter, r *http.Request, params handler.GetUsersParams) {
	cW.UserApi.GetUsers(w, r, applicationparameter.UserQueryParams{
		Limit:  *params.Limit,
		Offset: *params.Offset,
	})
}

func (cW *userControllerWrap) PostUsers(w http.ResponseWriter, r *http.Request) {
	cW.UserApi.PostUsers(w, r)
}

func (cW *userControllerWrap) PutUsersUserId(w http.ResponseWriter, r *http.Request, userID string) {
	cW.UserApi.PutUsersUserId(w, r, userID)
}

func (cW *userControllerWrap) GetUsersUserId(w http.ResponseWriter, r *http.Request, userID string) {
	cW.UserApi.GetUsersUserId(w, r, userID)
}

func (cW *userControllerWrap) DeleteUsersUserId(w http.ResponseWriter, r *http.Request, userID string) {
	cW.UserApi.DeleteUsersUserId(w, r, userID)
}
