package main

import (
	"corona-information-service/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	// Handle port assignment (either based on environment variable, or local override)
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Set up handler endpoints
	http.HandleFunc(handler.DEFAULT_PATH, handler.DefaultHandler)
	http.HandleFunc(handler.CASE_PATH, handler.CaseHandler)
	http.HandleFunc(handler.POLICY_PATH, handler.PolicyHandler)
	http.HandleFunc(handler.STATUS_PATH, handler.StatusHandler)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
