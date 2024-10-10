package platforms

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/gin-pagination/v2/pkg/pagination"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestGetPlatforms(t *testing.T) {
	tests := []struct {
		Name       string
		Page       int
		PageSize   int
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Response   string
	}{
		{
			"GetPlatforms - sql error - GetPlatforms",
			0,
			10,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT p.(.+) FROM platforms").WithArgs(10, 0).WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetPlatforms - sql error - CountPlatforms",
			0,
			2,
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "website", "country", "source", "notes", "comment", "privacy", "contacts_count", "articles_count", "projects_count", "platform_categories"}).
					AddRow(1, "test", "test", "test", "test", "test", "test", "test", 1, 1, 1, "test").
					AddRow(2, "test", "test", "test", "test", "test", "test", "test", 1, 1, 1, "test")
				mock.ExpectQuery("SELECT p.(.+) FROM platforms").WithArgs(2, 0).WillReturnRows(rows)

				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM platforms").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		}, {
			"GetPlatforms - 2 platforms from page 1",
			0,
			2,
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "website", "country", "source", "notes", "comment", "privacy", "contacts_count", "articles_count", "projects_count", "platform_categories"}).
					AddRow(1, "test", "test", "test", "test", "test", "test", "test", 1, 1, 1, "test").
					AddRow(2, "test", "test", "test", "test", "test", "test", "test", 1, 1, 1, "test")
				mock.ExpectQuery("SELECT p.(.+) FROM platforms").WithArgs(2, 0).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"count"}).AddRow(10)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM platforms").WillReturnRows(rows)
			},
			http.StatusOK,
			`{"total":10,"platforms":[{"id":1,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"name":"test","website":"test","country":"test","source":"test","notes":"test","privacy":"test","comment":"test","categories":null,"categoryString":"test","contactsCount":1,"articlesCount":1,"projectsCount":1},{"id":2,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"name":"test","website":"test","country":"test","source":"test","notes":"test","privacy":"test","comment":"test","categories":null,"categoryString":"test","contactsCount":1,"articlesCount":1,"projectsCount":1}]}`,
		},
		{
			"GetPlatforms - 4 platforms from page 2",
			1,
			4,
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "website", "country", "source", "notes", "comment", "privacy", "contacts_count", "articles_count", "projects_count", "platform_categories"}).
					AddRow(3, "test", "test", "test", "test", "test", "test", "test", 1, 1, 1, "test").
					AddRow(4, "test", "test", "test", "test", "test", "test", "test", 1, 1, 1, "test").
					AddRow(5, "test", "test", "test", "test", "test", "test", "test", 1, 1, 1, "test").
					AddRow(6, "test", "test", "test", "test", "test", "test", "test", 1, 1, 1, "test")
				mock.ExpectQuery("SELECT p.(.+) FROM platforms").WithArgs(4, 4).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"count"}).AddRow(10)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM platforms").WillReturnRows(rows)
			},
			http.StatusOK,
			`{"total":10,"platforms":[{"id":3,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"name":"test","website":"test","country":"test","source":"test","notes":"test","privacy":"test","comment":"test","categories":null,"categoryString":"test","contactsCount":1,"articlesCount":1,"projectsCount":1},{"id":4,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"name":"test","website":"test","country":"test","source":"test","notes":"test","privacy":"test","comment":"test","categories":null,"categoryString":"test","contactsCount":1,"articlesCount":1,"projectsCount":1},{"id":5,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"name":"test","website":"test","country":"test","source":"test","notes":"test","privacy":"test","comment":"test","categories":null,"categoryString":"test","contactsCount":1,"articlesCount":1,"projectsCount":1},{"id":6,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"name":"test","website":"test","country":"test","source":"test","notes":"test","privacy":"test","comment":"test","categories":null,"categoryString":"test","contactsCount":1,"articlesCount":1,"projectsCount":1}]}`,
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
			r.GET("/api/platforms",
				pagination.New(
					pagination.WithSizeText("pageSize"),
					pagination.WithMinPageSize(1),
					pagination.WithMaxPageSize(100),
				),
				GetPlatforms(env),
			)

			// Create httptest request
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/platforms?page=%d&pageSize=%d", test.Page, test.PageSize), nil)
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
