package middleware

import (
	"net/http"

	"oms/backend/api/v1/helpers"
)

// PanicRecoveryMiddleware recovers from panics and returns 500 error
func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// TODO: Log panic with stack trace
				// logger := logrus.WithField("panic", err)
				// logger.Error("Panic recovered")

				helpers.WriteErrorResponse(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

