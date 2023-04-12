package users

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	// This timestamp is to mock date values returned by the database
	timestamp, err := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
	if err != nil {
		t.Fatal(err)
	}

	sqlTimestamp := sql.NullTime{
		Time:  timestamp,
		Valid: false,
	}

	tests := []struct {
		Name       string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Body       string
		Response   string
	}{
		{
			"Login - missing email",
			nil,
			http.StatusBadRequest,
			`{"password": "test"}`,
			`{"error":"Key: 'LoginInput.Email' Error:Field validation for 'Email' failed on the 'required' tag"}`,
		},
		{
			"Login - missing password",
			nil,
			http.StatusBadRequest,
			`{"email": "test"}`,
			`{"error":"Key: 'LoginInput.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			"Login - sql error - GetUserWithEmail",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"email": "test", "password": "test"}`,
			`{}`,
		},
		{
			"Login - Incorrect Password",
			func(mock sqlmock.Sqlmock) {
				// Hash a fake password for testing
				hashBytes, err := bcrypt.GenerateFromPassword([]byte("testsomething"), bcrypt.DefaultCost)
				if err != nil {
					panic(err)
				}
				hashedPassword := string(hashBytes)

				rows := sqlmock.NewRows([]string{"id", "created_at", "modified_at", "deleted_at", "email", "password", "role"}).
					AddRow(1, timestamp, timestamp, sqlTimestamp, "test", hashedPassword, "user")
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

			},
			http.StatusUnauthorized,
			`{"email": "test", "password": "nottest"}`,
			`{}`,
		},
		{
			"Login - Successfull login but an error on the jwt creation",
			func(mock sqlmock.Sqlmock) {
				// Hash a fake password for testing
				hashBytes, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
				if err != nil {
					panic(err)
				}
				hashedPassword := string(hashBytes)

				rows := sqlmock.NewRows([]string{"id", "created_at", "modified_at", "deleted_at", "email", "password", "role"}).
					// Userid 0 will force the mock jwt service to return an error
					AddRow(0, timestamp, timestamp, sqlTimestamp, "test", hashedPassword, "user")
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

			},
			http.StatusInternalServerError,
			`{"email": "test", "password": "test"}`,
			`{}`,
		},
		{
			"Login - Successfull login",
			func(mock sqlmock.Sqlmock) {
				// Hash a fake password for testing
				hashBytes, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
				if err != nil {
					panic(err)
				}
				hashedPassword := string(hashBytes)

				rows := sqlmock.NewRows([]string{"id", "created_at", "modified_at", "deleted_at", "email", "password", "role"}).
					AddRow(1, timestamp, timestamp, sqlTimestamp, "test", hashedPassword, "user")
				mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

			},
			http.StatusOK,
			`{"email": "test", "password": "test"}`,
			`{"token":"usertoken","user":{"id":1,"createdAt":"2023-01-01T00:00:00Z","modifiedAt":"2023-01-01T00:00:00Z","deletedAt":{"Valid":false,"Time":"0001-01-01T00:00:00Z"},"role":"user","email":"test"}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Initilize test router, environemnt and mock database
			r, mockDb, mockSql, env, err := utils.SetupTestEnvironment(test.MockDbCall)
			// Close the mock database at the end of the test
			defer mockDb.Close()

			// Check for errors during setup
			require.NoError(t, err)

			// Register handler
			r.POST("/login", Login(env))

			// Create httptest request
			req, _ := http.NewRequest("POST", "/login", strings.NewReader(test.Body))
			w := httptest.NewRecorder()

			// Mock request
			r.ServeHTTP(w, req)

			// Read response data
			responseData, _ := io.ReadAll(w.Body)

			// Check response status
			require.Equal(t, test.StatusCode, w.Code)

			// Handle empty responses
			response := string(responseData)
			if response == "" {
				response = "{}"
			}

			// Check response body
			require.JSONEq(t, test.Response, response)

			// Check for any remaining expectations
			// we make sure that all expectations were met
			if err := mockSql.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}

}
