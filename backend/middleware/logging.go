package middleware

import (
	"net/http"
	"time"

	// TODO: Add logging library import
	// "github.com/sirupsen/logrus"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// TODO: Use structured logging (Logrus)
		// logger := logrus.WithFields(logrus.Fields{
		//     "method": r.Method,
		//     "path":   r.URL.Path,
		//     "ip":     r.RemoteAddr,
		// })
		// logger.Info("Request started")

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		// TODO: Log response
		// logger.WithField("duration", duration).Info("Request completed")
		_ = duration // Placeholder to avoid unused variable
	})
}

