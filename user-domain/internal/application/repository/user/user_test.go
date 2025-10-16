package userrepository

import (
	"context"
	"errors"
	"testing"
	application_mock "user-domain/internal/application/mocks/outbound"
	domainoutport "user-domain/internal/domain/outport"
	"user-domain/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newRepository(t *testing.T) (re domainoutport.UserRepository, m *application_mock.UserRepo) {
	repo := application_mock.NewUserRepo(t)
	return NewUserRepo(repo), repo
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	userID := uuid.New().String()
	tests := []struct {
		name    string
		data    *entity.User
		mock    func(repo *application_mock.UserRepo)
		wantErr bool
	}{
		{
			name: "success",
			data: &entity.User{
				ID:    userID,
				Name:  "user test",
				Email: "usertest@gmail.com",
				Phone: "0123456789",
			},
			mock: func(repo *application_mock.UserRepo) {
				u := &entity.User{
					ID:    userID,
					Name:  "user test",
					Email: "usertest@gmail.com",
					Phone: "0123456789",
				}
				repo.On("CreateUser", mock.Anything, u).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid id",
			data: &entity.User{
				ID:    "id",
				Name:  "user test",
				Email: "usertest@gmail.com",
				Phone: "0123456789",
			},
			mock: func(repo *application_mock.UserRepo) {
				u := &entity.User{
					ID:    "id",
					Name:  "user test",
					Email: "usertest@gmail.com",
					Phone: "0123456789",
				}
				repo.On("CreateUser", mock.Anything, u).Return(errors.New("invalid error"))
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			data: &entity.User{
				ID:    userID,
				Name:  "user test",
				Email: "usertest.com",
				Phone: "0123456789",
			},
			mock: func(repo *application_mock.UserRepo) {
				u := &entity.User{
					ID:    userID,
					Name:  "user test",
					Email: "usertest.com",
					Phone: "0123456789",
				}
				repo.On("CreateUser", mock.Anything, u).Return(errors.New("invalid error")).Maybe()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userRepo, repomock := newRepository(t)
			t.Parallel()
			tt.mock(repomock)
			err := userRepo.CreateUser(context.Background(), tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			repomock.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	userID := uuid.New()
	tests := []struct {
		name    string
		data    *entity.User
		mock    func(repo *application_mock.UserRepo)
		wantErr bool
	}{
		{
			name: "success",
			data: &entity.User{
				ID:   userID.String(),
				Name: "son edit",
			},
			mock: func(repo *application_mock.UserRepo) {
				data := &entity.User{
					ID:   userID.String(),
					Name: "son edit",
				}
				repo.On("UpdateUser", mock.Anything, data).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, repomock := newRepository(t)
			tt.mock(repomock)
			err := userRepo.UpdateUser(t.Context(), tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			repomock.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	validUserID := uuid.New().String()
	invalidUserID := "abd"
	t.Parallel()
	tests := []struct {
		name    string
		mock    func(m *application_mock.UserRepo)
		userID  string
		wantErr bool
	}{
		{
			name: "success",
			mock: func(m *application_mock.UserRepo) {
				m.On("DeleteUser", mock.Anything, validUserID).Return(nil)
			},
			userID: validUserID,
		},
		{
			name: "faild invalid user id",
			mock: func(m *application_mock.UserRepo) {
				m.On("DeleteUser", mock.Anything, invalidUserID).Return(nil)
			},
			userID: invalidUserID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userRepo, repomock := newRepository(t)
			tt.mock(repomock)
			err := userRepo.DeleteUser(t.Context(), tt.userID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
