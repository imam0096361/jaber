package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"news-portal/src/config"
	"news-portal/src/database"
	"news-portal/src/router"
)

func init() {
	// Change to project root if in src directory
	if wd, err := os.Getwd(); err == nil {
		if filepath.Base(wd) == "src" {
			os.Chdir("..")
		}
	}
}

func main() {
	// Initialize database
	database.Connect()
	defer database.Close()

	// Insert sample data
	database.InsertSampleData()

	// Create HTTP mux
	mux := http.NewServeMux()

	// Setup routes
	router.SetupRoutes(mux)

	// Start server
	port := config.SERVER_PORT
	fmt.Printf("News Portal running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
