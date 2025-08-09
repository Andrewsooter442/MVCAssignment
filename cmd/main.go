package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/Andrewsooter442/MVCAssignment/handler"
	"github.com/Andrewsooter442/MVCAssignment/internal/model"
	"github.com/Andrewsooter442/MVCAssignment/middleware"
)

func main() {
	// Set environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Creating a DB connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DBNAME"),
	)
	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pool := model.ModelConnection{DB: db}
	app := &handler.Application{Pool: pool}

	// Creating a Server and listen
	mux := http.NewServeMux()

	//Can add any middleware before them, that's why using(Handle and HandleFunc)
	//mux.Handle("/", middleware.VerifyJWTMiddleware(http.HandlerFunc(handler.HandleRootRequest)))
	// removing jwt verification
	mux.Handle("/", middleware.VerifyJWT(http.HandlerFunc(app.HandleRootRequest)))
	mux.Handle("/login", (http.HandlerFunc(app.HandleLoginRequest)))
	mux.Handle("/signup", (http.HandlerFunc(app.HandleSignupRequest)))
	mux.Handle("/api/", middleware.VerifyJWT(http.HandlerFunc(app.HandleApiRequest)))
	mux.Handle("/admin/", middleware.VerifyJWT(http.HandlerFunc(app.HandleAdminRequest)))

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), (mux)); err != nil {
		fmt.Println("Error starting server")
		log.Fatal(err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
