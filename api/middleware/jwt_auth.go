package middleware

import (
	"context"
	"main/internal"
	"net/http"
	"strings"
)

const (
	userIDKey string = "x-user-id"
)

//nolint:mnd // business logic
func JwtAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			t := strings.Split(authHeader, " ")
			if len(t) == 2 {
				authToken := t[1]
				authorized, err := internal.IsAuthorized(authToken, secret)
				if authorized {
					var userID string
					userID, err = internal.ExtractIDFromToken(authToken, secret)
					if err != nil {
						http.Error(w, jsonError(err.Error()), http.StatusUnauthorized)
						return
					}
					ctx := context.WithValue(r.Context(), userIDKey, userID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				http.Error(w, jsonError(err.Error()), http.StatusUnauthorized)
				return
			}
			http.Error(w, jsonError("Not authorized"), http.StatusUnauthorized)
		})
	}
}

func jsonError(message string) string {
	return `{"message": "` + message + `"}`
}
