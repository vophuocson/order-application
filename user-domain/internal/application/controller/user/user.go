package usercontroller

import (
	"encoding/json"
	"net/http"
	"user-domain/internal/application/controller/apiutil"
	applicationparameter "user-domain/internal/application/controller/parameter"
	"user-domain/internal/application/controller/user/dto"
	applicationerror "user-domain/internal/application/error"
	applicationinbound "user-domain/internal/application/inbound"
	applicationoutbound "user-domain/internal/application/outbound"
	domaininport "user-domain/internal/domain/inport"

	"user-domain/internal/entity"
)

type user struct {
	sv     domaininport.UserService
	logger applicationoutbound.Logger
}

func (h *user) PostUsers(w http.ResponseWriter, r *http.Request) {
	responseWriter := apiutil.NewJSONResponse(w, r, h.logger)
	userDto := dto.UserPost{}
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		responseWriter.Failure(apiutil.WrapError(err, applicationerror.ErrDecode))
		return
	}
	userEntity := entity.User{}
	userDto.MapTo(&userEntity)
	err = h.sv.CreateUser(r.Context(), &userEntity)
	if err != nil {
		responseWriter.Failure(err)
		return
	}
	responseWriter.Success(http.StatusCreated, nil)
}

func (h *user) PutUsersUserId(w http.ResponseWriter, r *http.Request, userID string) {
	responseWriter := apiutil.NewJSONResponse(w, r, h.logger)
	userDto := dto.UserPut{}
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		responseWriter.Failure(apiutil.WrapError(err, applicationerror.ErrDecode))
		return
	}
	userEntity := entity.User{ID: userID}
	userDto.MapTo(&userEntity)
	err = h.sv.UpdateUser(r.Context(), &userEntity)
	if err != nil {
		responseWriter.Failure(err)
	}
	responseWriter.Success(http.StatusNoContent, nil)
}

func (h *user) GetUsersUserId(w http.ResponseWriter, r *http.Request, userID string) {
	responseWriter := apiutil.NewJSONResponse(w, r, h.logger)
	userEntity, err := h.sv.GetUserByID(r.Context(), userID)
	if err != nil {
		responseWriter.Failure(err)
		return
	}
	res := dto.UserResponse{}
	res.GetFrom(userEntity)
	responseWriter.Success(http.StatusOK, res)
}

func (h *user) DeleteUsersUserId(w http.ResponseWriter, r *http.Request, userID string) {
	responseWriter := apiutil.NewJSONResponse(w, r, h.logger)
	err := h.sv.DeleteUser(r.Context(), userID)
	if err != nil {
		responseWriter.Failure(err)
	}
}

func (h *user) GetUsers(w http.ResponseWriter, r *http.Request, paramObj applicationparameter.UserQueryParams) {
	responseWriter := apiutil.NewJSONResponse(w, r, h.logger)
	eUsers, err := h.sv.ListUsers(r.Context(), paramObj.Offset, paramObj.Limit)
	if err != nil {
		responseWriter.Failure(err)
		return
	}
	users := dto.UsersResponse{}
	users.GetFrom(eUsers)
	responseWriter.Success(http.StatusOK, users)
}

func NewUserControler(sv domaininport.UserService, logger applicationoutbound.Logger) applicationinbound.UserApi {
	return &user{
		sv:     sv,
		logger: logger,
	}
}
