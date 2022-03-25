package main

import (
	"corona-information-service/internal/db"
	"corona-information-service/internal/handler"
	"corona-information-service/internal/model"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	//Initializes firestore client
	db.InitializeFirestore()

	// Handle port assignment (either based on environment variable, or local override)
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Set up handler endpoints
	http.HandleFunc(model.DEFAULT_PATH, handler.DefaultHandler)
	http.HandleFunc(model.CASE_PATH, handler.CaseHandler)
	http.HandleFunc(model.POLICY_PATH, handler.PolicyHandler)
	http.HandleFunc(model.STATUS_PATH, handler.StatusHandler)
	http.HandleFunc(model.NOTIFICATION_PATH, handler.NotificationHandler)
	http.HandleFunc(strings.TrimSuffix(model.NOTIFICATION_PATH, "/"), handler.NotificationHandler) //Will be forgiving since some forget "/" at the end

	// Start server
	log.Println("Starting server on port " + port + " ...")

	//Will allow firestore to close the connection before terminating the server
	err := http.ListenAndServe(":"+port, nil)
	db.CloseFirestore()
	log.Fatal(err)
}
