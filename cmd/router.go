package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Andrewsooter442/MVCAssignment/handler"
	"github.com/Andrewsooter442/MVCAssignment/middleware"
)

func registerRoutes(mainMux *http.ServeMux, app *handler.Application) {

	mainMux.Handle("/", middleware.VerifyJWT(http.HandlerFunc(app.HandleRootRequest)))
	mainMux.Handle("/login", (http.HandlerFunc(app.HandleLoginRequest)))
	mainMux.Handle("/signup", (http.HandlerFunc(app.HandleSignupRequest)))

	// Admin routes
	adminRouter := mux.NewRouter()
	adminSubrouter := adminRouter.PathPrefix("/admin").Subrouter()
	adminSubrouter.Use(middleware.VerifyJWT)
	adminSubrouter.Use(middleware.CheckAdmin)

	adminSubrouter.HandleFunc("/editItem/{id}", app.HandleGetEditItem).Methods("GET")
	adminSubrouter.HandleFunc("/editItem/{id}", app.HandlePostEditItem).Methods("POST")
	adminSubrouter.HandleFunc("/addItem", app.HandleGetAddItem).Methods("GET")
	adminSubrouter.HandleFunc("/addItem", app.HandlePostAddItem).Methods("POST")
	adminSubrouter.HandleFunc("/addCategory", app.HandleGetAddCategory).Methods("GET")
	adminSubrouter.HandleFunc("/addCategory", app.HandlePostAddCategory).Methods("POST")
	adminSubrouter.HandleFunc("/viewOldOrders", app.HandleGetViewOldOrder).Methods("GET")
	adminSubrouter.HandleFunc("/vieworder/{id}", app.HandleGetViewOrder).Methods("GET")

	mainMux.Handle("/admin/", adminRouter)

	// Api routes
	apiRouter := mux.NewRouter()
	apiSubrouter := apiRouter.PathPrefix("/api").Subrouter()
	apiSubrouter.Use(middleware.VerifyJWT)

	apiSubrouter.HandleFunc("/logout", app.HandleLogout).Methods("GET")
	apiSubrouter.HandleFunc("/placeOrder", app.HandlePlaceOrder).Methods("POST")
	apiSubrouter.HandleFunc("/payment", app.HandleGetPayment).Methods("Get")
	apiSubrouter.HandleFunc("/payment", app.HandlePostPayment).Methods("POST")
	apiSubrouter.HandleFunc("/completeOrder", app.HandleCompleteOrderItem).Methods("POST")

	mainMux.Handle("/api/", apiRouter)
}
