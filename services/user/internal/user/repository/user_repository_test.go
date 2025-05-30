package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/model"
	"github.com/alifmufthi91/ecommerce-system/services/user/internal/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, data model.User)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name string
		data model.User
		sqlMock
		wantErr bool
	}{
		{
			name: "success",
			data: model.User{
				Name:         "Test User",
				Phone:        "+6281234567890",
				Email:        "testuser@example.com",
				PasswordHash: "$2a$10$EIXo1z5Zb1f8Q3e5j5k6uO0d9F4h5l7m8n9o0p1q2r3s4t5u6v7w8x",
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.User) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "users" ("name","phone","email","password_hash","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
						),
					).WithArgs(
						data.Name,
						data.Phone,
						data.Email,
						data.PasswordHash,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).WillReturnRows(
						sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()),
					)
				},
			},
			wantErr: false,
		},
		{
			name: "error - failed to create user with existing phone",
			data: model.User{
				Name: "Test User",
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.User) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "users" ("name","phone","email","password_hash","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
						),
					).WithArgs(
						data.Name,
						data.Phone,
						data.Email,
						data.PasswordHash,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).WillReturnError(
						errors.New("users_phone_key already exists"),
					)
				},
			},
			wantErr: true,
		},
		{
			name: "error - failed to create user with existing email",
			data: model.User{
				Name:         "Test User",
				Email:        "testuser@example.com",
				Phone:        "+6281234567890",
				PasswordHash: "$2a$10$EIXo1z5Zb1f8Q3e5j5k6uO0d9F4h5l7m8n9o0p1q2r3s4t5u6v7w8x",
			},
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, data model.User) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(
							`INSERT INTO "users" ("name","phone","email","password_hash","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`,
						),
					).WithArgs(
						data.Name,
						data.Phone,
						data.Email,
						data.PasswordHash,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).WillReturnError(
						errors.New("users_email_key already exists"),
					)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.data)

			repo := NewUserRepository(mockDb.Db)

			err := repo.CreateUser(context.Background(), &tt.data)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}

func TestGetUserByEmailOrPhone(t *testing.T) {
	mockDb, err := pkg.SetupMockDB()

	type sqlMock struct {
		Setup func(mockDB sqlmock.Sqlmock, emailOrPhone string, user model.User)
	}

	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}
	tests := []struct {
		name         string
		emailOrPhone string
		sqlMock      sqlMock
		wantUser     model.User
		wantErr      bool
	}{
		{
			name:         "success - found by email",
			emailOrPhone: "testuser@example.com",
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, emailOrPhone string, user model.User) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 OR phone = $2 ORDER BY "users"."id" LIMIT $3`),
					).WithArgs(
						emailOrPhone,
						emailOrPhone,
						1,
					).WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"}).
							AddRow(uuid.New(), user.Name, user.Email, user.Phone, time.Now(), time.Now()),
					)
				},
			},
			wantUser: model.User{
				Name:  "Test User",
				Email: "testuser@example.com",
				Phone: "+6281234567890",
			},
			wantErr: false,
		},
		{
			name:         "success - found by phone",
			emailOrPhone: "+6281234567890",
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, emailOrPhone string, user model.User) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 OR phone = $2 ORDER BY "users"."id" LIMIT $3`),
					).WithArgs(
						emailOrPhone,
						emailOrPhone,
						1,
					).WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "email", "phone", "created_at", "updated_at"}).
							AddRow(uuid.New(), user.Name, user.Email, user.Phone, time.Now(), time.Now()),
					)
				},
			},
			wantUser: model.User{
				Name:  "Test User",
				Email: "testuser@example.com",
				Phone: "+6281234567890",
			},
			wantErr: false,
		},
		{
			name:         "error - user not found",
			emailOrPhone: "",
			sqlMock: sqlMock{
				Setup: func(mockDB sqlmock.Sqlmock, emailOrPhone string, user model.User) {
					mockDB.ExpectQuery(
						regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 OR phone = $2 ORDER BY "users"."id" LIMIT $3`),
					).WithArgs(
						emailOrPhone,
						emailOrPhone,
						1,
					).WillReturnError(
						gorm.ErrRecordNotFound,
					)
				},
			},
			wantUser: model.User{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.sqlMock.Setup(mockDb.Mock, tt.emailOrPhone, tt.wantUser)

			repo := NewUserRepository(mockDb.Db)

			user, err := repo.GetUserByEmailOrPhone(context.Background(), tt.emailOrPhone)

			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.wantUser.Name, user.Name)
			assert.Equal(t, tt.wantUser.Email, user.Email)
			assert.Equal(t, tt.wantUser.Phone, user.Phone)
		})
	}
}
