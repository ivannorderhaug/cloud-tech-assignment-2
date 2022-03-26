package main

import (
	"log"
	"net/http"
	"os"
)

const CLIENT_PATH = "/client"

func main() {
	// Handle port assignment (either based on environment variable, or local override)
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8081")
		port = "8081"
	}

	http.HandleFunc(CLIENT_PATH, nil)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
