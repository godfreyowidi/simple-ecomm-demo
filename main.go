package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/godfreyowidi/simple-ecomm-demo/db"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Postgres connection
	postgresDB, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer postgresDB.Close()

	// Example test route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Server is running and connected to DB")
	})

	// Get PORT from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("✅ Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
