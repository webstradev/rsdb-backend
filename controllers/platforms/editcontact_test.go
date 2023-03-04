package platforms

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestEditContact(t *testing.T) {
	tests := []struct {
		Name             string
		IdString         string
		PlatformIdString string
		MockDbCall       func(sqlmock.Sqlmock)
		StatusCode       int
		Body             string
		Response         string
	}{
		{
			"EditContact - non int id",
			"notanint",
			"notanint",
			nil,
			http.StatusBadRequest,
			`{}`,
			`{"error":"Invalid ID"}`,
		},
		{
			"EditContact - non int platformId",
			"1",
			"notanint",
			nil,
			http.StatusBadRequest,
			`{}`,
			`{"error":"Invalid Platform ID"}`,
		},
		{
			"EditContact - Bad json body",
			"1",
			"1",
			nil,
			http.StatusBadRequest,
			`{badbody}`,
			`{"error": "invalid character 'b' looking for beginning of object key string"}`,
		},
		{
			"EditContact - sql error on EditPlatform",
			"1",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE contacts SET").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{"name":"test","title":"test","email":"test","phone":"","phone2":"","address":"","notes":"","source":"test","privacy":"test"}`,
			`{}`,
		},
		{
			"EditPlatform - Valid Request",
			"1",
			"1",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE contacts SET").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			http.StatusOK,
			`{"name":"test","title":"test","email":"test","phone":"","phone2":"","address":"","notes":"","source":"test","privacy":"test"}`,
			`{}`,
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
			r.PUT("/api/v1/platforms/:platformId/contacts/:id", EditContact(env))

			// Create httptest request
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/platforms/%s/contacts/%s", test.PlatformIdString, test.IdString), strings.NewReader(test.Body))
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
