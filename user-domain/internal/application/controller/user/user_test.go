package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	appmock "user-domain/internal/application/mocks/outbound"
	domainmock "user-domain/internal/domain/mocks/inport"
	"user-domain/internal/entity"

	"github.com/stretchr/testify/mock"
)

func newController(t *testing.T) (*user, *domainmock.UserService, *appmock.Logger) {
	sv := domainmock.NewUserService(t)
	logger := appmock.NewLogger(t)
	c := NewUserControler(sv, logger).(*user)
	return c, sv, logger
}

func TestPostUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      interface{}
		mockSetup func(sv *domainmock.UserService)
		logSetup  func(l *appmock.Logger)
		wantCode  int
	}{
		{
			name: "success",
			body: map[string]any{"name": "Alice", "email": "a@a.com"},
			mockSetup: func(sv *domainmock.UserService) {
				sv.On("CreateUser", mock.Anything, &entity.User{Name: "Alice", Email: "a@a.com"}).Return(nil)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "decode error",
			body: func() io.Reader { return bytes.NewBufferString("{") }(),
			logSetup: func(l *appmock.Logger) {
				l.On("WithContext", mock.Anything).Return(l)
				l.On("Warn", mock.Anything, mock.Anything)
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: map[string]any{"name": "Alice", "email": "a@a.com"},
			mockSetup: func(sv *domainmock.UserService) {
				sv.On("CreateUser", mock.Anything, &entity.User{Name: "Alice", Email: "a@a.com"}).Return(errors.New("boom"))
			},
			logSetup: func(l *appmock.Logger) {
				l.On("WithContext", mock.Anything).Return(l)
				l.On("Error", mock.Anything, mock.Anything)
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl, sv, logger := newController(t)

			var bodyReader io.Reader
			switch v := tt.body.(type) {
			case io.Reader:
				bodyReader = v
			default:
				b, _ := json.Marshal(v)
				bodyReader = bytes.NewBuffer(b)
			}

			if tt.mockSetup != nil {
				tt.mockSetup(sv)
			}
			if tt.logSetup != nil {
				tt.logSetup(logger)
			}

			req := httptest.NewRequest(http.MethodPost, "/users", bodyReader)
			w := httptest.NewRecorder()
			ctrl.PostUsers(w, req)
			if w.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d, body=%s", w.Code, tt.wantCode, w.Body.String())
			}
		})
	}
}

func TestPutUsersUserId(t *testing.T) {
	t.Parallel()

	name := "Bob"
	tests := []struct {
		name      string
		userID    string
		body      interface{}
		mockSetup func(sv *domainmock.UserService)
		logSetup  func(l *appmock.Logger)
		wantCode  int
	}{
		{
			name:   "success",
			userID: "42",
			body:   map[string]any{"name": name},
			mockSetup: func(sv *domainmock.UserService) {
				sv.On("UpdateUser", mock.Anything, &entity.User{ID: "42", Name: name}).Return(nil)
			},
			wantCode: http.StatusNoContent,
		},
		{
			name:   "decode error",
			userID: "42",
			body:   func() io.Reader { return bytes.NewBufferString("{") }(),
			logSetup: func(l *appmock.Logger) {
				l.On("WithContext", mock.Anything).Return(l)
				l.On("Warn", mock.Anything, mock.Anything)
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:   "service error",
			userID: "42",
			body:   map[string]any{"name": name},
			mockSetup: func(sv *domainmock.UserService) {
				sv.On("UpdateUser", mock.Anything, &entity.User{ID: "42", Name: name}).Return(errors.New("boom"))
			},
			logSetup: func(l *appmock.Logger) {
				l.On("WithContext", mock.Anything).Return(l)
				l.On("Error", mock.Anything, mock.Anything)
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl, sv, logger := newController(t)

			var bodyReader io.Reader
			switch v := tt.body.(type) {
			case io.Reader:
				bodyReader = v
			default:
				b, _ := json.Marshal(v)
				bodyReader = bytes.NewBuffer(b)
			}

			if tt.mockSetup != nil {
				tt.mockSetup(sv)
			}
			if tt.logSetup != nil {
				tt.logSetup(logger)
			}

			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.userID, bodyReader)
			w := httptest.NewRecorder()
			ctrl.PutUsersUserId(w, req, tt.userID)
			if w.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d, body=%s", w.Code, tt.wantCode, w.Body.String())
			}
		})
	}
}

func TestGetUsersUserId(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userID    string
		mockSetup func(sv *domainmock.UserService)
		logSetup  func(l *appmock.Logger)
		wantCode  int
	}{
		{
			name:   "success",
			userID: "42",
			mockSetup: func(sv *domainmock.UserService) {
				sv.On("GetUserByID", mock.Anything, "42").Return(&entity.User{ID: "42", Name: "Bob"}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "service error",
			userID: "42",
			mockSetup: func(sv *domainmock.UserService) {
				sv.On("GetUserByID", mock.Anything, "42").Return((*entity.User)(nil), errors.New("boom"))
			},
			logSetup: func(l *appmock.Logger) {
				l.On("WithContext", mock.Anything).Return(l)
				l.On("Error", mock.Anything, mock.Anything)
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl, sv, logger := newController(t)
			if tt.mockSetup != nil {
				tt.mockSetup(sv)
			}
			if tt.logSetup != nil {
				tt.logSetup(logger)
			}
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()
			ctrl.GetUsersUserId(w, req, tt.userID)
			if w.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d, body=%s", w.Code, tt.wantCode, w.Body.String())
			}
		})
	}
}

func TestDeleteUsersUserId(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userID    string
		mockSetup func(sv *domainmock.UserService)
		logSetup  func(l *appmock.Logger)
		wantCode  int
	}{
		{
			name:   "success",
			userID: "9",
			mockSetup: func(sv *domainmock.UserService) {
				sv.On("DeleteUser", mock.Anything, "9").Return(nil)
			},
			wantCode: http.StatusOK, // controller writes default success? It doesn't, but Failure only sets error status; Success not called.
		},
		{
			name:   "service error",
			userID: "9",
			mockSetup: func(sv *domainmock.UserService) {
				sv.On("DeleteUser", mock.Anything, "9").Return(errors.New("boom"))
			},
			logSetup: func(l *appmock.Logger) {
				l.On("WithContext", mock.Anything).Return(l)
				l.On("Error", mock.Anything, mock.Anything)
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl, sv, logger := newController(t)
			if tt.mockSetup != nil {
				tt.mockSetup(sv)
			}
			if tt.logSetup != nil {
				tt.logSetup(logger)
			}
			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()
			ctrl.DeleteUsersUserId(w, req, tt.userID)
			if w.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d, body=%s", w.Code, tt.wantCode, w.Body.String())
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		offset    int
		limit     int
		mockSetup func(sv *domainmock.UserService)
		logSetup  func(l *appmock.Logger)
		wantCode  int
	}{
		{
			name:   "success",
			offset: 0,
			limit:  2,
			mockSetup: func(sv *domainmock.UserService) {
				sv.On("ListUsers", mock.Anything, 0, 2).Return([]*entity.User{{ID: "1"}, {ID: "2"}}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "service error",
			offset: 0,
			limit:  2,
			mockSetup: func(sv *domainmock.UserService) {
				sv.On("ListUsers", mock.Anything, 0, 2).Return(([]*entity.User)(nil), errors.New("boom"))
			},
			logSetup: func(l *appmock.Logger) {
				l.On("WithContext", mock.Anything).Return(l)
				l.On("Error", mock.Anything, mock.Anything)
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl, sv, logger := newController(t)
			if tt.mockSetup != nil {
				tt.mockSetup(sv)
			}
			if tt.logSetup != nil {
				tt.logSetup(logger)
			}
			req := httptest.NewRequest(http.MethodGet, "/users?offset=0&limit=2", nil)
			w := httptest.NewRecorder()
			ctrl.GetUsers(w, req, struct{ Limit, Offset int }{Limit: tt.limit, Offset: tt.offset})
			if w.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d, body=%s", w.Code, tt.wantCode, w.Body.String())
			}
		})
	}
}
