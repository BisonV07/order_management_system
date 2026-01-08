package middleware

import (
	"context"
	"net/http"
	"strings"

	"oms/server/api/v1/helpers"
	"oms/server/core/auth"
)

// AuthMiddleware extracts and validates JWT token from Authorization header
// Sets user_id in request context for downstream handlers
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public endpoints
		publicPaths := []string{
			"/api/v1/health",
			"/api/v1/products",
			"/api/v1/auth/login",
			"/api/v1/auth/signup",
		}
		for _, path := range publicPaths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}
		// Also allow GET requests to products with ID
		if r.Method == "GET" && len(r.URL.Path) > 15 && r.URL.Path[:15] == "/api/v1/products" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helpers.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "Missing Authorization header")
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "Invalid Authorization header format")
			return
		}

		token := parts[1]

		// Validate JWT token and extract user_id and role
		userID, role, err := auth.ValidateToken(token)
		if err != nil {
			helpers.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized", "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "user_role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

