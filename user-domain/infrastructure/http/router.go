package http

import (
	"net/http"
	"user-domain/infrastructure/http/handler"
	"user-domain/infrastructure/http/middleware"
	postgresuser "user-domain/infrastructure/persistence/postgres/user"
	"user-domain/internal/application/controller/parameter"
	controlleruser "user-domain/internal/application/controller/user"
	"user-domain/internal/application/inbound"
	"user-domain/internal/application/logger"
	"user-domain/internal/application/orchestrator"
	"user-domain/internal/application/outbound"
	repositoryuser "user-domain/internal/application/repository/user"
	domainuser "user-domain/internal/domain/user"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func BuildRouter(db *gorm.DB, logger outbound.Logger, workflowRuner outbound.WorkflowRuner) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.LoggingMiddleware(logger))
	r.Route("/api/v1", func(r chi.Router) {
		buildUserSubRouter(r, db, logger, workflowRuner)
	})
	return r
}

func buildUserSubRouter(r chi.Router, db *gorm.DB, loggerOutbound outbound.Logger, workflowRuner outbound.WorkflowRuner) {
	userPersistence := postgresuser.NewUserRepo(db)
	userRepo := repositoryuser.NewUserRepo(userPersistence)

	loggerOutport := logger.NewLogger(loggerOutbound)

	var producer outbound.Producer
	var subscriber outbound.Subscriber
	o := orchestrator.NewWorkflowStarter(producer, subscriber, loggerOutbound, workflowRuner)

	userService := domainuser.NewUserService(userRepo, loggerOutport, o)

	userControler := controlleruser.NewUserControler(userService, loggerOutbound)

	handler.HandlerWithOptions(&userControllerWrap{UserApi: userControler}, handler.ChiServerOptions{
		BaseRouter: r,
	})
}

type userControllerWrap struct {
	inbound.UserApi
}

func (cW *userControllerWrap) GetUsers(w http.ResponseWriter, r *http.Request, params handler.GetUsersParams) {
	cW.UserApi.GetUsers(w, r, parameter.UserQueryParams{
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
