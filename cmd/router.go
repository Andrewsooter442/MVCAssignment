package main

import (
	"net/http"

	"github.com/Andrewsooter442/MVCAssignment/handler"
	"github.com/Andrewsooter442/MVCAssignment/middleware"
)

func registerRoutes(mux *http.ServeMux, app *handler.Application) {

	//Can add any middleware before them, that's why using(Handle and HandleFunc)
	//mux.Handle("/", middleware.VerifyJWTMiddleware(http.HandlerFunc(handler.HandleRootRequest)))
	// removing jwt verification
	mux.Handle("/", middleware.VerifyJWT(http.HandlerFunc(app.HandleRootRequest)))
	mux.Handle("/login", (http.HandlerFunc(app.HandleLoginRequest)))
	mux.Handle("/signup", (http.HandlerFunc(app.HandleSignupRequest)))
	mux.Handle("/api/", middleware.VerifyJWT(http.HandlerFunc(app.HandleApiRequest)))
	mux.Handle("/admin/", middleware.VerifyJWT(http.HandlerFunc(app.HandleAdminRequest)))

}
