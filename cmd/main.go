package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Andrewsooter442/MVCAssignment/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	log.Println("Running database migrations...")
	driver, err := mysql.WithInstance(app.Pool.DB, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalf("migration setup failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("an error occurred while running migrations: %v", err)
	}
	log.Println("Migrations completed successfully.")

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
