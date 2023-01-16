package ghostmon_test

import (
	"github.com/justpark/ghostmon"
	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/require"
	"io"
	"net"
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
			_, w := net.Pipe()
			request := httptest.NewRequest(testCase.method, "/unpostpone", nil)
			response := ServeAndHandleRequest(t, w, request)

			require.Equal(t, http.StatusMethodNotAllowed, response.Code)
		})
	}
}

func TestHandleUnpostponeSendsCommand(t *testing.T) {
	r, w := net.Pipe()
	request := httptest.NewRequest(http.MethodPost, "/unpostpone", nil)

	var wg conc.WaitGroup
	wg.Go(func() {
		response := ServeAndHandleRequest(t, w, request)
		require.Equal(t, http.StatusCreated, response.Code)
	})

	b, _ := io.ReadAll(r)
	require.Equal(t, "unpostpone", string(b))

	wg.Wait()
}

func ServeAndHandleRequest(
	t *testing.T,
	writer net.Conn,
	request *http.Request,
) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	communicator := ghostmon.NewCommunicator(writer)
	server := ghostmon.NewHTTPServer(communicator)
	server.Handler.ServeHTTP(recorder, request)

	return recorder
}
