package articles

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/webstradev/gin-pagination/v2/pkg/pagination"
	"github.com/webstradev/rsdb-backend/utils"
)

func TestGetArticles(t *testing.T) {
	// This timestamp is to mock date values returned by the database
	timestamp, err := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name       string
		Page       int
		PageSize   int
		MockDbCall func(sqlmock.Sqlmock)
		StatusCode int
		Response   string
	}{
		{
			"GetArticles - sql error - GetArticles",
			0,
			10,
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT a.(.+) FROM articles").WithArgs(10, 0).WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetArticles - sql error - CountArticles",
			0,
			2,
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "link", "date", "body"}).
					AddRow(1, "test", "test", "test", sql.NullTime{Valid: true, Time: timestamp}, "test").
					AddRow(2, "test", "test", "test", sql.NullTime{Valid: true, Time: timestamp}, "test")
				mock.ExpectQuery("SELECT a.(.+) FROM articles").WithArgs(2, 0).WillReturnRows(rows)

				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM articles").WillReturnError(errors.New("test"))
			},
			http.StatusInternalServerError,
			`{}`,
		},
		{
			"GetArticles - 2 articles from page 1",
			0,
			2,
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "link", "date", "body"}).
					AddRow(1, "test", "test", "test", sql.NullTime{Valid: true, Time: timestamp}, "test").
					AddRow(2, "test", "test", "test", sql.NullTime{Valid: true, Time: timestamp}, "test")
				mock.ExpectQuery("SELECT a.(.+) FROM articles").WithArgs(2, 0).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"id"}).AddRow(10)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM articles").WillReturnRows(rows)
			},
			http.StatusOK,
			`{"total":10,"articles":[{"id":1,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"title":"test","description":"test","link":"test","date":{"Time":"2023-01-01T00:00:00Z","Valid":true},"body":"test","platforms":null,"tags":null,"tagString":""},{"id":2,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"title":"test","description":"test","link":"test","date":{"Time":"2023-01-01T00:00:00Z","Valid":true},"body":"test","platforms":null,"tags":null,"tagString":""}]}`,
		},
		{
			"GetArticles - 4 articles from page 2",
			1,
			4,
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "link", "date", "body"}).
					AddRow(3, "test", "test", "test", sql.NullTime{Valid: true, Time: timestamp}, "test").
					AddRow(4, "test", "test", "test", sql.NullTime{Valid: true, Time: timestamp}, "test").
					AddRow(5, "test", "test", "test", sql.NullTime{Valid: true, Time: timestamp}, "test").
					AddRow(6, "test", "test", "test", sql.NullTime{Valid: true, Time: timestamp}, "test")
				mock.ExpectQuery("SELECT a.(.+) FROM articles").WithArgs(4, 4).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"id"}).AddRow(10)
				mock.ExpectQuery("SELECT COUNT(.+) AS count FROM articles").WillReturnRows(rows)
			},
			http.StatusOK,
			`{"total":10,"articles":[{"id":3,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"title":"test","description":"test","link":"test","date":{"Time":"2023-01-01T00:00:00Z","Valid":true},"body":"test","platforms":null,"tags":null,"tagString":""},{"id":4,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"title":"test","description":"test","link":"test","date":{"Time":"2023-01-01T00:00:00Z","Valid":true},"body":"test","platforms":null,"tags":null,"tagString":""},{"id":5,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"title":"test","description":"test","link":"test","date":{"Time":"2023-01-01T00:00:00Z","Valid":true},"body":"test","platforms":null,"tags":null,"tagString":""},{"id":6,"createdAt":"0001-01-01T00:00:00Z","modifiedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false},"title":"test","description":"test","link":"test","date":{"Time":"2023-01-01T00:00:00Z","Valid":true},"body":"test","platforms":null,"tags":null,"tagString":""}]}`,
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
			r.GET("/api/articles", pagination.New(pagination.WithSizeText("pageSize"), pagination.WithMinPageSize(1), pagination.WithMaxPageSize(100)), GetArticles(env))

			// Create httptest request
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/articles?page=%d&pageSize=%d", test.Page, test.PageSize), nil)
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
