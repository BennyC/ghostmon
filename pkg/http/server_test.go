package http_test

import (
	gtest "github.com/justpark/ghostmon/pkg/testing"
	"github.com/sourcegraph/conc"
	"github.com/steinfletcher/apitest-jsonpath/jsonpath"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
	httpServer, listener := gtest.CreateHTTPServer(t)

	var wg conc.WaitGroup
	wg.Go(func() {
		msg := make([]byte, 6)
		_, _ = listener.Read(msg)
		w, _ := os.ReadFile("fixtures/status.txt")
		_, _ = listener.Write(w)

		require.EqualValues(t, "status", msg)
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/status", nil)
	httpServer.Handler.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.NoError(t, jsonpath.Equal("table", "`test`.`sample_data_0`", response.Body))

	wg.Wait()
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
