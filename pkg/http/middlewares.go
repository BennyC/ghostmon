package http

import "net/http"

// OnlyAllowMethod will perform checks on the http.Request Method and check if it
// equals the allowed method. When a disallowed method is presented, a 405 Status
// will be returned to the client
func OnlyAllowMethod(allowedMethod string, next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != allowedMethod {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		next(writer, request)
	}
}
