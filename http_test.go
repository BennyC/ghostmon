package ghostmon_test

import (
	"bytes"
	"github.com/justpark/ghostmon"
	"github.com/stretchr/testify/require"
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
			request := httptest.NewRequest(testCase.method, "/unpostpone", nil)
			writer := httptest.NewRecorder()
			serveRequest(t, writer, request)

			require.Equal(t, http.StatusMethodNotAllowed, writer.Code)
		})
	}
}

func serveRequest(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
	t.Helper()

	var b bytes.Buffer
	communicator := ghostmon.NewCommunicator(&b)
	server := ghostmon.NewHTTPServer(communicator)
	server.Handler.ServeHTTP(recorder, request)
}
