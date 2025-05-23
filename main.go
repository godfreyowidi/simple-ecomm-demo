package main

import (
	"log"
	"net/http"
	"os"

	"github.com/godfreyowidi/simple-ecomm-demo/db"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Init DB
	if err := db.InitDB(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Example test route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running and connected to DB"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
