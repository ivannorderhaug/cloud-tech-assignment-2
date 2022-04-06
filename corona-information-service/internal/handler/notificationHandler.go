package handler

import (
	"corona-information-service/tools/customjson"
	"corona-information-service/tools/webhook"
	"net/http"
	"strings"
)

// NotificationHandler */
func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		response, err := webhook.RegisterWebhook(r)
		if err != nil {
			http.Error(w, "Error in registering webhook", http.StatusInternalServerError)
			return
		}
		customjson.Encode(w, response)
	}

	if r.Method == http.MethodGet || r.Method == http.MethodDelete {
		methodHandler(w, r)
	}
}

// methodHandler */
func methodHandler(w http.ResponseWriter, r *http.Request) {
	//Splits url into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	switch len(parts) {
	//If length of parts is 4, then it means the user has wants all webhooks
	case 4:
		if r.Method != http.MethodGet {
			http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
			return
		}
		encodeAllWebhooks(w)

	//If the length of parts is 5, then it means the user has specified webhook id in their request
	case 5:
		if r.Method == http.MethodDelete {
			encodeWebhookDeletionResponse(w, parts[4])
		}

		if r.Method == http.MethodGet {
			encodeSingleWebhook(w, parts[4])
		}
	default:
		http.Error(w, "Incorrect path format.", http.StatusBadRequest)
		return
	}
}

// encodeAllWebhooks */
func encodeAllWebhooks(w http.ResponseWriter) {
	webhooks, err := webhook.GetAllWebhooks()
	if err != nil {
		http.Error(w, "Error retrieving webhooks from database", http.StatusInternalServerError)
		return
	}

	if len(webhooks) == 0 {
		http.Error(w, "There are currently no webhooks registered in the database", http.StatusNotFound)
		return
	}

	customjson.Encode(w, webhooks)
}

// encodeSingleWebhook */
func encodeSingleWebhook(w http.ResponseWriter, id string) {
	wh, found := webhook.GetWebhook(id)
	if !found {
		http.Error(w, "Unable to locate webhook in database", http.StatusNotFound)
		return
	}
	customjson.Encode(w, wh)
}

// encodeWebhookDeletionResponse */
func encodeWebhookDeletionResponse(w http.ResponseWriter, id string) {
	deleted, err := webhook.DeleteWebhook(id)
	if err != nil || !deleted {
		http.Error(w, "Error removing webhook from database. it might not exist", http.StatusInternalServerError)
		return
	}
	response := make(map[string]string, 1)
	response["result"] = "The webhook has been successfully removed from the database!"
	customjson.Encode(w, response)
}
