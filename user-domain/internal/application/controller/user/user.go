package userapplication

import (
	"encoding/json"
	"net/http"
	"user-domain/controller/dto"
	"user-domain/controller/handler"
	"user-domain/controller/inbound"
	userapioutbound "user-domain/internal/application/controller/user/outbound"
	userinport "user-domain/internal/domain/user/inport"

	"user-domain/internal/entity"
)

type user struct {
	sv     userinport.UserService
	logger userapioutbound.Logger
}

func (h *user) PostUsers(w http.ResponseWriter, r *http.Request) {
	userDtoRequest := handler.PostUsersJSONRequestBody{}
	err := json.NewDecoder(r.Body).Decode(&userDtoRequest)
	if err != nil {
		// handle error here
	}
	var userDto = createUserPostFromPostUsersRequestObject(&userDtoRequest)
	userEntity := entity.User{}
	userDto.MapTo(&userEntity)
	err = h.sv.CreateUser(r.Context(), &userEntity)
	if err != nil {
		// handle error here
	}
}

func createUserPostFromPostUsersRequestObject(s *handler.PostUsersJSONRequestBody) *dto.UserPost {
	return &dto.UserPost{
		Name:  s.Name,
		Email: string(s.Email),
		// Phone:   s.Phone,
		// Address: s.Address,
	}
}

// func (api *user) Create(userReq *dto.UserPost) error {
// 	userEntity := entity.User{}
// 	userReq.MapTo(&userEntity)
// 	ctx := context.Background()
// 	err := api.sv.CreateUser(ctx, &userEntity)
// 	return err
// }

// func (api *user) Update(userID string, userReq *dto.UserPut) error {
// 	userEntity := entity.User{ID: userID}
// 	userReq.MapTo(&userEntity)
// 	ctx := context.Background()
// 	err := api.sv.UpdateUser(ctx, &userEntity)
// 	return err
// }

// func (api *user) Get(userID string) (*dto.UserResponse, error) {
// 	ctx := context.Background()
// 	userEntity, err := api.sv.GetUserByID(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res := dto.UserResponse{}
// 	res.GetFrom(userEntity)
// 	return &res, nil
// }

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

func NewUserControler(sv userinport.UserService, logger userapioutbound.Logger) inbound.UserApi {
	return &user{
		sv:     sv,
		logger: logger,
	}
}
