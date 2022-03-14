package handler

import "net/http"

// DefaultHandler /*
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "No functionality on this level. Please use "+CASE_PATH+", "+POLICY_PATH+", "+STATUS_PATH+" or "+NOTIFICATION_PATH+".\nYou can also get more information from the README. ", http.StatusOK)
}
