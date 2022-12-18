package middlewares

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestPaginationMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		queryParams      url.Values
		expectedPage     int
		expectedPageSize int
	}{
		{
			"Non int Page Param - Bad Request",
			url.Values{
				"page": {"notanumber"},
			},
			0,
			0,
		},
		{
			"Non int PageSize Param - Bad Request",
			url.Values{
				"page":     {"1"},
				"pageSize": {"notanumber"},
			},
			0,
			0,
		},
		{
			"Default Handling",
			url.Values{},
			1,
			10,
		},
		{
			"The first 100 results",
			url.Values{
				"page":     {"1"},
				"pageSize": {"100"},
			},
			1,
			100,
		},
		{
			"The second 20 results",
			url.Values{
				"page":     {"2"},
				"pageSize": {"20"},
			},
			2,
			20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test context
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			// Add the query parameters to the Request of the test context
			ctx.Request = &http.Request{
				URL: &url.URL{
					RawQuery: url.Values(tt.queryParams).Encode(),
				},
			}

			// Call middleware on the test context
			PaginationMiddleware()(ctx)

			// Make assertions on expected page and pageSize
			require.Equal(t, tt.expectedPage, ctx.GetInt("page"))
			require.Equal(t, tt.expectedPageSize, ctx.GetInt("pageSize"))
		})
	}
}
