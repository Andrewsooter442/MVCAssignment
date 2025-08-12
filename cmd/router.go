package main

import (
	"net/http"

	"github.com/Andrewsooter442/MVCAssignment/handler"
	"github.com/Andrewsooter442/MVCAssignment/middleware"
	"github.com/gorilla/mux"
)

func registerRoutes(mainMux *http.ServeMux, app *handler.Application) {

	mainMux.Handle("/", middleware.VerifyJWT(http.HandlerFunc(app.HandleRootRequest)))
	mainMux.Handle("/login", (http.HandlerFunc(app.HandleLoginRequest)))
	mainMux.Handle("/signup", (http.HandlerFunc(app.HandleSignupRequest)))
	mainMux.Handle("/admin/", middleware.VerifyJWT(http.HandlerFunc(app.HandleAdminRequest)))

	apiRouter := mux.NewRouter()

	apiSubrouter := apiRouter.PathPrefix("/api").Subrouter()
	apiSubrouter.Use(middleware.VerifyJWT) // Apply your JWT middleware here

	apiSubrouter.HandleFunc("/editItem/{id}", app.HandleGetEditItem).Methods("GET")
	apiSubrouter.HandleFunc("/editItem/{id}", app.HandlePostEditItem).Methods("POST")
	//apiSubrouter.HandleFunc("/products", app.HandleCreateProduct).Methods("POST")
	//apiSubrouter.HandleFunc("/products/{sku}", app.HandleUpdateProduct).Methods("PUT")
	apiSubrouter.HandleFunc("/logout", app.HandleLogout).Methods("GET")

	mainMux.Handle("/api/", apiRouter)
}
