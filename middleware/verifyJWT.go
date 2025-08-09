package middleware

import (
	"context"
	"fmt"
	"github.com/Andrewsooter442/MVCAssignment/config"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		fmt.Println(cookie, err.Error())
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value
		claims := &config.JWTtoken{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), config.UserObject, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
