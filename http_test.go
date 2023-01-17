package ghostmon_test

import (
	"github.com/justpark/ghostmon"
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
			request := httptest.NewRequest(testCase.method, "/unpostpone", nil)
			response := httptest.NewRecorder()
			httpServer, _, closer := CreateServer(t)
			httpServer.Handler.ServeHTTP(response, request)
			defer closer()

			require.Equal(t, http.StatusMethodNotAllowed, response.Code)
		})
	}
}

func TestHandleUnpostponeSendsCommand(t *testing.T) {
	httpServer, ghostServer, closer := CreateServer(t)
	defer closer()

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/unpostpone", nil)
	httpServer.Handler.ServeHTTP(response, request)

	require.Equal(t, http.StatusCreated, response.Code)
	b, _ := io.ReadAll(ghostServer)

	require.Equal(t, "unpostpone", string(b))
}

func CreateServer(t *testing.T) (*http.Server, net.Conn, func()) {
	t.Helper()

	// TODO Handle err
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	var ghostServer net.Conn
	go func() {
		// TODO Handle err
		ghostServer, _ = ln.Accept()
	}()

	communicator, err := ghostmon.NewNetCommunicator(ln.Addr())
	if err != nil {
		t.FailNow()
	}

	return ghostmon.NewHTTPServer(communicator), ghostServer, func() {
		// TODO Handle err
		_ = ln.Close()
	}
}
