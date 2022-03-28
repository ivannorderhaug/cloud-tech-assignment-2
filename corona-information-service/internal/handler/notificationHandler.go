package handler

import (
	"corona-information-service/tools"
	"net/http"
	"strings"
)

// NotificationHandler */
func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		response, err := tools.RegisterWebhook(r)
		if err != nil {
			http.Error(w, "Error in registering webhook", http.StatusInternalServerError)
			return
		}
		tools.Encode(w, response)
	}

	if r.Method == http.MethodGet || r.Method == http.MethodDelete {
		pathHandler(w, r)
	}
}

// pathHandler */
func pathHandler(w http.ResponseWriter, r *http.Request) {
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
	webhooks, err := tools.GetAllWebhooks()
	if err != nil {
		http.Error(w, "Error retrieving webhooks from database", http.StatusInternalServerError)
		return
	}

	if len(webhooks) == 0 {
		http.Error(w, "There are currently no webhooks registered in the database", http.StatusNotFound)
		return
	}

	tools.Encode(w, webhooks)
}

// encodeSingleWebhook */
func encodeSingleWebhook(w http.ResponseWriter, id string) {
	webhook, found := tools.GetWebhook(id)
	if !found {
		http.Error(w, "Unable to locate webhook in database", http.StatusNotFound)
		return
	}
	tools.Encode(w, webhook)
}

// encodeWebhookDeletionResponse */
func encodeWebhookDeletionResponse(w http.ResponseWriter, id string) {
	deleted, err := tools.DeleteWebhook(id)
	if err != nil || !deleted {
		http.Error(w, "Error removing webhook from database. it might not exist", http.StatusInternalServerError)
		return
	}
	response := make(map[string]string, 1)
	response["result"] = "The webhook has been successfully removed from the database!"
	tools.Encode(w, response)
}
