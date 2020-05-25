package app

import (
	"net/http"
	u "qt-api/utils"
	"strings"
	"qt-api/models"
	jwt "github.com/dgrijalva/jwt-go"
	"os"
	"context"
	"time"
	// "fmt"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		noAuth := []string{"/v1/users/create", "/v1/users/authenticate", "/v1/teststest"} // Endpoints w/out auth
		requestPath := r.URL.Path

		// Serve requests w/out auth
		for _, value := range noAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")

		// Return '403 Unauthorized' if token is missing
		if tokenHeader == "" {
			u.Error(w, 403, "Forbidden")
			return
		}

		// Make sure token comes in the valid way
		splitStr := strings.Split(tokenHeader, " ")
		if len(splitStr) != 2 {
			u.Error(w, http.StatusForbidden, "Forbidden")
			return
		}

		authToken := splitStr[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(authToken, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		// If there is an error, the token was bad. Return 403
		if err != nil {
			u.Error(w, 403, "Forbidden")
			return
		}

		// Invalid for some other reason, probably not signed on backend
		if !token.Valid {
			u.Error(w, 403, "Forbidden")
			return
		}

		if time.Unix(tk.StandardClaims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
			u.Error(w, 401, "Unauthorized")
			return
		}
	

		// Token valid and all is good, continue
		ctx := context.WithValue(r.Context(), "userId", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	});
}