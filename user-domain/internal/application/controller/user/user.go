package usercontroller

import (
	"encoding/json"
	"errors"
	"net/http"
	"user-domain/internal/application/controller/apiutil"
	"user-domain/internal/application/controller/user/dto"
	applicationinbound "user-domain/internal/application/inbound"
	applicationoutbound "user-domain/internal/application/outbound"
	domainerror "user-domain/internal/domain/error"
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
		h.logger.WithContext(r.Context()).Warn("error decode: %s", err)
		responseWriter.Failure(err)
		return
	}
	userEntity := entity.User{}
	userDto.MapTo(&userEntity)
	err = h.sv.CreateUser(r.Context(), &userEntity)
	if err != nil {
		if errors.Is(err, domainerror.ErrCodeInternal) {
			h.logger.WithContext(r.Context()).Error("error decode: %s", err)
		} else {
			h.logger.WithContext(r.Context()).Warn("error decode: %s", err)
		}
		responseWriter.Failure(err)
		return
	}
	responseWriter.Success(http.StatusCreated, nil)
}

// func (api *user) Update(userID string, userReq *dto.UserPut) error {
// 	userEntity := entity.User{ID: userID}
// 	userReq.MapTo(&userEntity)
// 	ctx := context.Background()
// 	err := api.sv.UpdateUser(ctx, &userEntity)
// 	return err
// }

func (api *user) GetUsersUserId(w http.ResponseWriter, r *http.Request, userID string) {
	responseWriter := apiutil.NewJSONResponse(w, r, api.logger)
	userEntity, err := api.sv.GetUserByID(r.Context(), userID)
	if err != nil {
		responseWriter.Failure(err)
		return
	}
	res := dto.UserResponse{}
	res.GetFrom(userEntity)
	responseWriter.Success(http.StatusOK, res)
}

// func (api *user) List(offset int, limit int) (*dto.UsersResponse, error) {
// 	ctx := context.Background()
// 	eUsers, err := api.sv.ListUsers(ctx, offset, limit)
// 	if err != nil {
// 		return nil, err
// 	}
// 	users := dto.UsersResponse{}
// 	users.GetFrom(eUsers)
// 	return &users, nil
// }

// func (api *user) Delete(userID string) error {
// 	ctx := context.Background()
// 	err := api.sv.DeleteUser(ctx, userID)
// 	return err
// }

func NewUserControler(sv domaininport.UserService, logger applicationoutbound.Logger) applicationinbound.UserApi {
	return &user{
		sv:     sv,
		logger: logger,
	}
}
