package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		next.ServeHTTP(w, r)
	})

}
