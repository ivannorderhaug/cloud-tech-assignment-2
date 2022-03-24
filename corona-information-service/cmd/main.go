package main

import (
	"corona-information-service/internal/db"
	handler2 "corona-information-service/internal/handler"
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
	http.HandleFunc(model.DEFAULT_PATH, handler2.DefaultHandler)
	http.HandleFunc(model.CASE_PATH, handler2.CaseHandler)
	http.HandleFunc(model.POLICY_PATH, handler2.PolicyHandler)
	http.HandleFunc(model.STATUS_PATH, handler2.StatusHandler)
	http.HandleFunc(model.NOTIFICATION_PATH, handler2.NotificationHandler)
	http.HandleFunc(strings.TrimSuffix(model.NOTIFICATION_PATH, "/"), handler2.NotificationHandler) //Will be forgiving since some forget "/" at the end

	// Start server
	log.Println("Starting server on port " + port + " ...")

	//Will allow firestore to close the connection before terminating the server
	err := http.ListenAndServe(":"+port, nil)
	db.CloseFirestore()
	log.Fatal(err)
}
