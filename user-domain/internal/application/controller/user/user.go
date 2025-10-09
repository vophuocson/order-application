package usercontroller

import (
	"encoding/json"
	"net/http"
	"user-domain/internal/application/controller/user/dto"
	applicationinbound "user-domain/internal/application/inbound"
	applicationoutbound "user-domain/internal/application/outbound"
	domaininport "user-domain/internal/domain/inport"
	domainoutport "user-domain/internal/domain/outport"

	"user-domain/internal/entity"
)

type user struct {
	sv     domaininport.UserService
	logger applicationoutbound.Logger
}

func (h *user) PostUsers(w http.ResponseWriter, r *http.Request) {
	userDto := dto.UserPost{}
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		h.logger.WithContext(r.Context()).Warnf(err.Error())
		return
	}
	userEntity := entity.User{}
	userDto.MapTo(&userEntity)
	err = h.sv.CreateUser(r.Context(), &userEntity)
	if err != nil {
		h.logger.WithContext(r.Context()).Warnf(err.Error())
		return
	}

}

// func (api *user) Update(userID string, userReq *dto.UserPut) error {
// 	userEntity := entity.User{ID: userID}
// 	userReq.MapTo(&userEntity)
// 	ctx := context.Background()
// 	err := api.sv.UpdateUser(ctx, &userEntity)
// 	return err
// }

func (api *user) Get(w http.ResponseWriter, r *http.Request) {
	var userID = r.PathValue("user_id")

	userEntity, err := api.sv.GetUserByID(r.Context(), userID)
	if err != nil {
		api.logger.WithContext(r.Context()).Error("error log", domainoutport.LogFields{})
		return
	}
	res := dto.UserResponse{}
	res.GetFrom(userEntity)
	entityEncoding, err := json.Marshal(&res)
	if err != nil {
		api.logger.WithContext(r.Context()).Error("err", domainoutport.LogFields{})
	}
	w.Write(entityEncoding)
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
