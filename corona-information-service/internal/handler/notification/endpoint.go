package _notification

import (
	"corona-information-service/tools/customjson"
	"corona-information-service/tools/webhook"
	"net/http"
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
