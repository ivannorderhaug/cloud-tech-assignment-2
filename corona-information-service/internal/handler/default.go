package handler

import (
	"corona-information-service/internal/model"
	"net/http"
)

// DefaultHandler */
func DefaultHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "No functionality on this level. Please use "+model.CASE_PATH+", "+model.POLICY_PATH+", "+model.STATUS_PATH+" or "+model.NOTIFICATION_PATH+".\nYou can also get more information from the README. ", http.StatusOK)
	}
}
