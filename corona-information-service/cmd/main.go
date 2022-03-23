package main

import (
	"corona-information-service/db"
	"corona-information-service/handler"
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
	http.HandleFunc(handler.DEFAULT_PATH, handler.DefaultHandler)
	http.HandleFunc(handler.CASE_PATH, handler.CaseHandler)
	http.HandleFunc(handler.POLICY_PATH, handler.PolicyHandler)
	http.HandleFunc(handler.STATUS_PATH, handler.StatusHandler)
	http.HandleFunc(handler.NOTIFICATION_PATH, handler.NotificationHandler)
	http.HandleFunc(strings.TrimSuffix(handler.NOTIFICATION_PATH, "/"), handler.NotificationHandler) //Will be forgiving since some forget "/" at the end

	// Start server
	log.Println("Starting server on port " + port + " ...")

	//Will allow firestore to close the connection before terminating the server
	err := http.ListenAndServe(":"+port, nil)
	db.CloseFirestore()
	log.Fatal(err)
}
