package userdomain

import (
    "context"
    "errors"
    "testing"

    domaininport "user-domain/internal/domain/inport"
    domainmock "user-domain/internal/domain/mocks/outport"
    "user-domain/internal/entity"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func newSvc(t *testing.T) (context.Context, *domainmock.UserRepository, *domainmock.Logger, domaininport.UserService) {
    ctx := context.Background()
    repoMock := domainmock.NewUserRepository(t)
    loggerMock := domainmock.NewLogger(t)
    svc := NewUserService(repoMock, loggerMock)
    return ctx, repoMock, loggerMock, svc
}

func TestCreateUser(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name       string
        setupMock  func(r *domainmock.UserRepository)
        input      *entity.User
        wantErr    error
    }{
        {
            name: "success",
            setupMock: func(r *domainmock.UserRepository) {
                r.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
            },
            input:   &entity.User{ID: "1", Name: "Alice"},
            wantErr: nil,
        },
        {
            name: "error from repo",
            setupMock: func(r *domainmock.UserRepository) {
                r.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("create failed"))
            },
            input:   &entity.User{ID: "1", Name: "Alice"},
            wantErr: errors.New("create failed"),
        },
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            ctx, repoMock, _, svc := newSvc(t)
            if tt.setupMock != nil {
                tt.setupMock(repoMock)
            }
            err := svc.CreateUser(ctx, tt.input)
            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.EqualError(t, err, tt.wantErr.Error())
                return
            }
            assert.NoError(t, err)
        })
    }
}

func TestGetUserByID(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name       string
        id         string
        setupMock  func(r *domainmock.UserRepository)
        want       *entity.User
        wantErr    error
    }{
        {
            name: "success",
            id:   "42",
            setupMock: func(r *domainmock.UserRepository) {
                r.On("GetUserByID", mock.Anything, "42").Return(&entity.User{ID: "42", Name: "Bob"}, nil)
            },
            want: &entity.User{ID: "42", Name: "Bob"},
        },
        {
            name: "error from repo",
            id:   "404",
            setupMock: func(r *domainmock.UserRepository) {
                r.On("GetUserByID", mock.Anything, "404").Return((*entity.User)(nil), errors.New("not found"))
            },
            want:    nil,
            wantErr: errors.New("not found"),
        },
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            ctx, repoMock, _, svc := newSvc(t)
            if tt.setupMock != nil {
                tt.setupMock(repoMock)
            }
            got, err := svc.GetUserByID(ctx, tt.id)
            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Nil(t, got)
                assert.EqualError(t, err, tt.wantErr.Error())
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}

func TestUpdateUser(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name       string
        setupMock  func(r *domainmock.UserRepository)
        input      *entity.User
        wantErr    error
    }{
        {
            name: "success",
            setupMock: func(r *domainmock.UserRepository) {
                r.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)
            },
            input:   &entity.User{ID: "7", Name: "Eve"},
            wantErr: nil,
        },
        {
            name: "error from repo",
            setupMock: func(r *domainmock.UserRepository) {
                r.On("UpdateUser", mock.Anything, mock.Anything).Return(errors.New("update failed"))
            },
            input:   &entity.User{ID: "7", Name: "Eve"},
            wantErr: errors.New("update failed"),
        },
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            ctx, repoMock, _, svc := newSvc(t)
            if tt.setupMock != nil {
                tt.setupMock(repoMock)
            }
            err := svc.UpdateUser(ctx, tt.input)
            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.EqualError(t, err, tt.wantErr.Error())
                return
            }
            assert.NoError(t, err)
        })
    }
}

func TestDeleteUser(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name       string
        id         string
        setupMock  func(r *domainmock.UserRepository)
        wantErr    error
    }{
        {
            name: "success",
            id:   "9",
            setupMock: func(r *domainmock.UserRepository) {
                r.On("DeleteUser", mock.Anything, "9").Return(nil)
            },
        },
        {
            name: "error from repo",
            id:   "9",
            setupMock: func(r *domainmock.UserRepository) {
                r.On("DeleteUser", mock.Anything, "9").Return(errors.New("delete failed"))
            },
            wantErr: errors.New("delete failed"),
        },
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            ctx, repoMock, _, svc := newSvc(t)
            if tt.setupMock != nil {
                tt.setupMock(repoMock)
            }
            err := svc.DeleteUser(ctx, tt.id)
            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.EqualError(t, err, tt.wantErr.Error())
                return
            }
            assert.NoError(t, err)
        })
    }
}

func TestListUsers(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name       string
        offset     int
        limit      int
        setupMock  func(r *domainmock.UserRepository)
        want       []*entity.User
        wantErr    error
    }{
        {
            name:   "success",
            offset: 0,
            limit:  10,
            setupMock: func(r *domainmock.UserRepository) {
                r.On("ListUsers", mock.Anything, 0, 10).Return([]*entity.User{{ID: "1", Name: "Alice"}, {ID: "2", Name: "Bob"}}, nil)
            },
            want: []*entity.User{{ID: "1", Name: "Alice"}, {ID: "2", Name: "Bob"}},
        },
        {
            name:   "error from repo",
            offset: 5,
            limit:  5,
            setupMock: func(r *domainmock.UserRepository) {
                r.On("ListUsers", mock.Anything, 5, 5).Return(([]*entity.User)(nil), errors.New("list failed"))
            },
            want:    nil,
            wantErr: errors.New("list failed"),
        },
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            ctx, repoMock, _, svc := newSvc(t)
            if tt.setupMock != nil {
                tt.setupMock(repoMock)
            }
            got, err := svc.ListUsers(ctx, tt.offset, tt.limit)
            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Nil(t, got)
                assert.EqualError(t, err, tt.wantErr.Error())
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}


