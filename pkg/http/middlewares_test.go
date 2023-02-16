package http_test

import (
	"testing"

	"net/http"
	"net/http/httptest"

	ghttp "github.com/justpark/ghostmon/pkg/http"
	"github.com/stretchr/testify/require"
)

func TestHandleOptionsMiddleware(t *testing.T) {
	handler := ghttp.OnlyAllowMethod(
		"POST",
		func(w http.ResponseWriter, r *http.Request) {

		},
	)

	response := httptest.NewRecorder()
	request := httptest.NewRequest("OPTIONS", "/", nil)
	handler(response, request)

	require.Equal(t, http.StatusOK, response.Code)
}
