package handler

import (
	"corona-information-service/internal/db"
	"corona-information-service/internal/model"
	"corona-information-service/tools"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const COLLECTION = "notifications"

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		registerWebhook(w, r)
		return
	}

	if r.Method == http.MethodGet || r.Method == http.MethodDelete {
		getWebhookHandler(w, r)
		return
	}

	http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
}

func getWebhookHandler(w http.ResponseWriter, r *http.Request) {
	//Splits url into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	switch len(parts) {
	//If length of parts is 4, then it means the user has wants all webhooks
	case 4:
		getAllWebhooks(w)
		return
	//If the length of parts is 5, then it means the user has specified webhook id in their request
	case 5:
		if r.Method == http.MethodDelete {
			deleteWebhook(w, parts[4])
		}

		if r.Method == http.MethodGet {
			getWebhook(w, parts[4])
		}
	default:
		http.Error(w, "Incorrect path format.", http.StatusBadRequest)
		return
	}
}

func getWebhook(w http.ResponseWriter, webhookId string) {
	documentFromFirestore, err := db.GetSingleDocumentFromFirestore(COLLECTION, webhookId)
	if err != nil {
		http.Error(w, "Could not find webhook with id: "+webhookId, http.StatusNotFound)
		return
	}

	var webhook model.Webhook
	webhook.ID = documentFromFirestore.Ref.ID
	documentFromFirestore.DataTo(&webhook)

	tools.Encode(w, webhook)
}

// getAllWebhooks */
func getAllWebhooks(w http.ResponseWriter) {
	documentsFromFirestore, err := db.GetAllDocumentsFromFirestore(COLLECTION)
	if err != nil {
		http.Error(w, "Error retrieving data from database", http.StatusInternalServerError)
		return
	}
	response := make([]model.Webhook, 0)

	//Converts each document snapshot into a webhook interface and adds it to the response slice
	for _, documentSnapshot := range documentsFromFirestore {
		var webhook model.Webhook
		webhook.ID = documentSnapshot.Ref.ID
		err = documentSnapshot.DataTo(&webhook)
		if err != nil {
			http.Error(w, "Error mapping data from database to data structure", http.StatusInternalServerError)
			return
		}
		response = append(response, webhook)
	}

	//If response slice is still 0 after loop that populates it, then it means there are no webhooks
	if len(response) == 0 {
		http.Error(w, "No webhooks found", http.StatusNotFound)
		return
	}

	tools.Encode(w, response)
}

// registerWebhook */
func registerWebhook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var wh model.Webhook

	err := decoder.Decode(&wh)
	if err != nil {
		http.Error(w, "Error during decoding", http.StatusInternalServerError)
		return
	}

	//checks if alpha3 code was used as param for country
	if len(wh.Country) == 3 {
		country, err := tools.GetCountryByAlphaCode(wh.Country)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		wh.Country = fmt.Sprint(country)
	}

	webhookID, err := db.AddToFirestore(COLLECTION, wh)
	if err != nil {
		http.Error(w, "Error adding webhook to database", http.StatusInternalServerError)
		return
	}

	//Respond with ID
	var response = make(map[string]string, 1)
	response["id"] = webhookID

	tools.Encode(w, response)
}

func deleteWebhook(w http.ResponseWriter, webhookId string) {
	if err := db.DeleteSingleDocumentFromFirestore(COLLECTION, webhookId); err != nil {
		http.Error(w, "Error deleting webhook from database!", http.StatusInternalServerError)
		return
	}

	response := make(map[string]string, 1)
	response["result"] = "The webhook has been successfully removed from the database!"

	tools.Encode(w, response)
}
