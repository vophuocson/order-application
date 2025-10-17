package userpersistence_test

import (
	"errors"
	"regexp"
	"testing"
	"time"
	userpersistence "user-domain/infrastructure/persistence/postgres/user"
	applicationoutbound "user-domain/internal/application/outbound"
	"user-domain/internal/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newNewUserRepo() (applicationoutbound.UserRepo, sqlmock.Sqlmock, error) {
	db, sqlmock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	g, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	if err != nil {
		return nil, nil, err
	}
	return userpersistence.NewUserRepo(g), sqlmock, nil
}

func TestCreateUser(t *testing.T) {
	invalidEmail := "invalidmail.com"
	repo, mock, err := newNewUserRepo()
	if err != nil {
		t.Error(err)
	}
	t.Parallel()
	tests := []struct {
		name    string
		mock    func(m sqlmock.Sqlmock)
		wantErr bool
		data    *entity.User
	}{
		{
			name: "success",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("id","deleted_at","email","phone","name") VALUES ($1,$2,$3,$4,$5) RETURNING "created_at","updated_at"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
						AddRow(time.Now(), time.Now()))
				m.ExpectCommit()
			},
			data: &entity.User{
				ID:    uuid.NewString(),
				Name:  "test",
				Email: "test@gmail.com",
				Phone: "12345678987654",
			},
			wantErr: false,
		},
		{
			name: "invalid mail",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("id","deleted_at","email","phone","name") VALUES ($1,$2,$3,$4,$5) RETURNING "created_at","updated_at"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("invalid email"))
				m.ExpectRollback()
			},
			data: &entity.User{
				Name:  "test",
				Email: invalidEmail,
				Phone: "12345678987654",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mock)
			err := repo.CreateUser(t.Context(), tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
