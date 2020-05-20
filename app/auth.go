package app

import (
	"net/http"
	u "qt-api/utils"
	"strings"
	"qt-api/models"
	jwt "github.com/dgrijalva/jwt-go"
	"os"
	"context"
	"fmt"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		noAuth := []string{} // Endpoints w/out auth
		requestPath := r.URL.Path

		// Serve requests w/out auth
		for _, value := range noAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string] interface{})
		tokenHeader := r.Header.Get("Authorization")

		// Return '403 Unauthorized' if token is missing
		if tokenHeader == "" {
			response = u.Message(false, "Missing token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		// Make sure token comes in the valid way
		splitStr := strings.Split(tokenHeader, " ")
		if len(splitStr) != 2 {
			response = u.Message(false, "Invalid/Malformed token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		authToken := splitStr[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(authToken, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		// If there is an error, the token was bad. Return 403
		if err != nil {
			response = u.Message(false, "Malformed token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		// Invalid for some other reason, probably not signed on backend
		if !token.Valid {
			response = u.Message(false, "Token not valid")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		// Token valid and all is good, log and continue
		fmt.Sprintf("User %", tk.UserId)
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	});
}