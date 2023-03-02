package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestGetContacts(t *testing.T) {
	tests := []struct {
		Name       string
		IdString   string
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Response   string
	}{
		{
			"GetContacts - non int id",
			"notanint",
			nil,
			http.StatusBadRequest,
			`{}`,
		},
		{
			"GetContacts - sql error on GetPlatform",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT .+ FROM contacts").WithArgs(1).WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetContacts - Valid Request",
			"1",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "title", "email", "phone", "phone2", "address", "notes", "source", "privacy", "platform_id"}).
					AddRow(1, "test", "test", "test", "test", "test", "test", "test", "test", "test", 1)
				mock.ExpectQuery("SELECT .+ FROM contacts").WithArgs(1).WillReturnRows(rows)
			},
			http.StatusOK,
			`[{"platformId":1,"id":1,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"name":"test","title":"test","email":"test","phone":"test","phone2":"test","address":"test","notes":"test","source":"test","privacy":"test"}]`,
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
			r.GET("/api/platforms/:id/contacts", GetContacts(env))

			// Create httptest request
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/platforms/%s/contacts", test.IdString), nil)
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
