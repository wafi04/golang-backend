package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/wafi04/golang-backend/grpc/pb"
	"github.com/wafi04/golang-backend/services/auth/repository"
	"golang.org/x/crypto/bcrypt"
)

func SetupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, *repository.UserRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewUserRepository(sqlxDB)

	return sqlxDB, mock, repo
}

func TestCreateUser(t *testing.T) {
	db, mock, repo := SetupMockDB(t)
	defer db.Close()

	req := &pb.CreateUserRequest{
		Name:       "Wafi",
		Email:      "wafiq610@gmail.com",
		Password:   "password123",
		IpAddress:  "127.0.0.1",
		DeviceInfo: "Android",
	}

	// Mock the user insertion
	mock.ExpectExec(`INSERT INTO users \(user_id, name, email, password_hash, role, is_active, is_email_verified, created_at, updated_at\)`).
		WithArgs(
			sqlmock.AnyArg(), // user_id
			req.Name,
			req.Email,
			req.Password,
			"admin",    // role (since email is wafiq610@gmail.com)
			true,       // is_active
			false,      // is_email_verified
			time.Now(), // created_at
			time.Now(), // updated_at
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO sessions`).
		WithArgs(
			sqlmock.AnyArg(), // session_id
			sqlmock.AnyArg(), // user_id
			sqlmock.AnyArg(), // access_token
			sqlmock.AnyArg(), // refresh_token
			req.IpAddress,
			req.DeviceInfo,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // last_activity_at
			true,             // is_active
			sqlmock.AnyArg(), // expires_at
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute the function
	resp, err := repo.CreateUser(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.UserId)
	assert.Equal(t, req.Name, resp.Name)
	assert.Equal(t, req.Email, resp.Email)
	assert.Equal(t, "admin", resp.Role)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.SessionInfo.SessionId)
	assert.Equal(t, req.DeviceInfo, resp.SessionInfo.DeviceInfo)
	assert.Equal(t, req.IpAddress, resp.SessionInfo.IpAddress)

	// Verify that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestLogin(t *testing.T) {
	db, mock, repo := SetupMockDB(t)
	defer db.Close()

	// Data dummy
	userID := "user-123"
	email := "wafiq610@gmail.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	rows := sqlmock.NewRows([]string{
		"user_id", "name", "email", "role", "password_hash", "picture", "is_email_verified", "created_at", "updated_at", "last_login_at", "is_active",
	}).AddRow(
		userID, "Wafi", email, "admin", string(hashedPassword), "", true, time.Now().Unix(), time.Now().Unix(), time.Now().Unix(), true,
	)
	mock.ExpectQuery(`SELECT`).
		WithArgs(email).
		WillReturnRows(rows)

	mock.ExpectQuery(`SELECT`).
		WithArgs(userID, "Android").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec(`INSERT INTO sessions`).
		WithArgs(
			sqlmock.AnyArg(),
			userID,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			"127.0.0.1",      // ip_address
			"Android",        // device_info
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // last_activity_at
			true,             // is_active
			sqlmock.AnyArg(), // expires_at
		).
		WillReturnError(fmt.Errorf("duplicate key")) // Simulasikan INSERT gagal

	// Mock query UPDATE sessions (akan dijalankan setelah INSERT gagal)
	mock.ExpectExec(`UPDATE sessions`).
		WithArgs(
			sqlmock.AnyArg(), // access_token
			sqlmock.AnyArg(), // refresh_token
			"127.0.0.1",      // ip_address
			sqlmock.AnyArg(), // last_activity_at
			userID,           // user_id
			"Android",        // device_info
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`UPDATE users`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	loginReq := &pb.LoginRequest{
		Email:      email,
		Password:   password,
		IpAddress:  "127.0.0.1",
		DeviceInfo: "Android",
	}
	resp, err := repo.Login(context.Background(), loginReq)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, userID, resp.UserId)
	assert.NotEmpty(t, resp.AccessToken)

	assert.NoError(t, mock.ExpectationsWereMet())
}
