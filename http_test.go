package ghostmon_test

import (
	"bytes"
	"github.com/justpark/ghostmon"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleUnpostponeOnlyAllowsPost(t *testing.T) {
	testCases := []struct {
		name   string
		method string
	}{
		{
			name:   "GET",
			method: http.MethodGet,
		},

		{
			name:   "PATCH",
			method: http.MethodPatch,
		},

		{
			name:   "DELETE",
			method: http.MethodDelete,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var b bytes.Buffer
			request := httptest.NewRequest(testCase.method, "/unpostpone", nil)
			response := serveRequest(t, &b, request)

			require.Equal(t, http.StatusMethodNotAllowed, response.Code)
		})
	}
}

func TestHandleUnpostponeSendsCommand(t *testing.T) {
	var b bytes.Buffer
	request := httptest.NewRequest(http.MethodPost, "/unpostpone", nil)
	response := serveRequest(t, &b, request)

	require.Equal(t, "unpostpone", b.String())
	require.Equal(t, http.StatusCreated, response.Code)
}

func serveRequest(
	t *testing.T,
	writer io.Writer,
	request *http.Request,
) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	communicator := ghostmon.NewCommunicator(writer)
	server := ghostmon.NewHTTPServer(communicator)
	server.Handler.ServeHTTP(recorder, request)

	return recorder
}
