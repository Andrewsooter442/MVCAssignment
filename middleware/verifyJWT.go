package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Andrewsooter442/MVCAssignment/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

const userObject = "user"

func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		fmt.Println(cookie, err.Error())
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value
		claims := &model.JWTtoken{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), userObject, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
