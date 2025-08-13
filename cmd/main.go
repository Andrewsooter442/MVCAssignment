package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Andrewsooter442/MVCAssignment/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load the html templates
	config.LoadTemplates()

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
	defer app.Pool.DB.Close()

	// Creating a Server and listen
	mainMux := http.NewServeMux()

	// Serve static files
	ServeStaticFiles(mainMux)

	//Register the routes
	registerRoutes(mainMux, app)

	//Start listning
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), (mainMux)); err != nil {
		fmt.Println("Error starting server")
		log.Fatal(err)
	}

}
