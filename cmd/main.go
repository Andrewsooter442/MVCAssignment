package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Set environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Creating a DB connection
	app, err := createDbConnection()
	if err != nil {
		log.Fatal(err)
	}

	// Creating a Server and listen
	mux := http.NewServeMux()

	//Register the routes
	registerRoutes(mux, app)

	//Start listning
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), (mux)); err != nil {
		fmt.Println("Error starting server")
		log.Fatal(err)
	}

}
