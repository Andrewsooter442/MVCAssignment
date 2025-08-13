package middleware

import (
	"github.com/Andrewsooter442/MVCAssignment/config"
	"net/http"
)

func CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(config.UserObject).(*config.JWTtoken)
		if !ok {
			http.Error(w, "Could not retrieve user claims", http.StatusInternalServerError)
			return
		}

		if claims.IsAdmin {
			next.ServeHTTP(w, r)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	})
}
