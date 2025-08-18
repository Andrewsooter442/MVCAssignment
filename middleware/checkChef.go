package middleware

import (
	"fmt"
	"github.com/Andrewsooter442/MVCAssignment/types"
	"net/http"
)

func CheckChef(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(types.UserObject).(*types.JWTtoken)
		if !ok {
			http.Error(w, "Could not retrieve user claims", http.StatusInternalServerError)
			return
		}

		if claims.IsCheff {
			next.ServeHTTP(w, r)
			return
		}
		fmt.Println("redirecting")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	})
}
