package ghostmon_test

import (
	gtest "github.com/justpark/ghostmon/pkg/testing"
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
			httpServer, listener := gtest.CreateServer(t)
			defer func() {
				_ = listener.Close()
			}()

			request := httptest.NewRequest(testCase.method, "/unpostpone", nil)
			response := httptest.NewRecorder()
			httpServer.Handler.ServeHTTP(response, request)

			require.Equal(t, http.StatusMethodNotAllowed, response.Code)
			// TODO Test we receive no connections to our server
		})
	}
}

func TestHandleUnpostponeSendsCommand(t *testing.T) {
	httpServer, listener := gtest.CreateServer(t)
	defer func() {
		_ = listener.Close()
	}()

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/unpostpone", nil)
	httpServer.Handler.ServeHTTP(response, request)
	msg, _ := gtest.ReceiveNetMessage(t, listener)

	require.Equal(t, http.StatusCreated, response.Code)
	require.EqualValues(t, "unpostpone", msg)
}
