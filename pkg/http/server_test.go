package http_test

import (
	gtest "github.com/justpark/ghostmon/pkg/testing"
	"github.com/sourcegraph/conc"
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
			httpServer, _ := gtest.CreateHTTPServer(t)
			request := httptest.NewRequest(testCase.method, "/unpostpone", nil)
			response := httptest.NewRecorder()

			httpServer.Handler.ServeHTTP(response, request)

			require.Equal(t, http.StatusMethodNotAllowed, response.Code)
		})
	}
}

func TestHandleUnpostponeSendsCommand(t *testing.T) {
	httpServer, listener := gtest.CreateHTTPServer(t)

	// TODO Figure out how to remove WaitGroups
	var wg conc.WaitGroup
	wg.Go(func() {
		msg, _ := io.ReadAll(listener)
		require.EqualValues(t, "unpostpone", msg)
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/unpostpone", nil)
	httpServer.Handler.ServeHTTP(response, request)

	require.Equal(t, http.StatusCreated, response.Code)
	wg.Wait()
}

func TestHandleStatusSendsStatusCommand(t *testing.T) {
	t.Skip()
}

func TestAbortSendsPanicCommand(t *testing.T) {
	httpServer, listener := gtest.CreateHTTPServer(t)

	var wg conc.WaitGroup
	wg.Go(func() {
		msg, _ := io.ReadAll(listener)
		require.EqualValues(t, "panic", msg)
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/abort", nil)
	httpServer.Handler.ServeHTTP(response, request)

	require.Equal(t, http.StatusCreated, response.Code)
	wg.Wait()
}
