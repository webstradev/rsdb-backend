package utils

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/webstradev/rsdb-backend/auth"
	"github.com/webstradev/rsdb-backend/db"
	"github.com/webstradev/rsdb-backend/mocks"
)

type Environment struct {
	DB          *db.Database
	JWT         auth.JWTServicer
	UUID        auth.UUIDGenerator
	AuthService auth.AuthServicer
}

func SetupTestEnvironment(MockDbCall func(sqlmock.Sqlmock)) (*gin.Engine, *sql.DB, sqlmock.Sqlmock, *Environment, error) {
	// Setup environment
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	env := Environment{}
	// Create mock database
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Create sqlx instance which an be passed to functions
	mockedDatabase := sqlx.NewDb(mockDb, "sqlmock")

	if MockDbCall != nil {
		MockDbCall(mockSql)
	}

	env.DB = db.SetupMockDB(mockedDatabase)

	// Create mock JWT service
	env.JWT = mocks.CreateMockJWTService()

	return r, mockDb, mockSql, &env, nil
}
