package tools

import (
	"bytes"
	"corona-information-service/internal/db"
	"corona-information-service/internal/model"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const COLLECTION = "notifications"

var webhooks = make([]model.Webhook, 0)

// InitializeWebhooks */
func InitializeWebhooks() {
	_, err := GetAllWebhooks()
	if err != nil {
		return
	}
}

// GetWebhook */
func GetWebhook(webhookId string) (model.Webhook, bool) {
	for _, wh := range webhooks {
		if wh.ID == webhookId {
			return wh, true
		}
	}
	return model.Webhook{}, false
}

// GetAllWebhooks */
func GetAllWebhooks() ([]model.Webhook, error) {
	if len(webhooks) == 0 {
		documentsFromFirestore, err := db.GetAllDocumentsFromFirestore(Hash([]byte(COLLECTION)))
		if err != nil {
			return []model.Webhook{}, err
		}

		//Converts each document snapshot into a webhook interface and adds it to the global webhooks slice
		for _, documentSnapshot := range documentsFromFirestore {
			var webhook model.Webhook
			err = documentSnapshot.DataTo(&webhook)
			if err != nil {
				return []model.Webhook{}, err
			}
			webhooks = append(webhooks, webhook)
		}
	}
	return webhooks, nil
}

// DeleteWebhook */
func DeleteWebhook(webhookId string) error {
	if len(webhooks) != 0 {
		for i, wh := range webhooks {
			if wh.ID == webhookId {
				webhooks = RemoveIndex(webhooks, i)
			}
		}
	}

	if err := db.DeleteSingleDocumentFromFirestore(Hash([]byte(COLLECTION)), Hash([]byte(webhookId))); err != nil {
		return err
	}

	return nil
}

// RegisterWebhook */
func RegisterWebhook(r *http.Request) (map[string]string, error) {
	var webhook model.Webhook

	err := Decode(r, &webhook)
	if err != nil {
		return map[string]string{}, err
	}

	//checks if alpha3 code was used as param for country
	if len(webhook.Country) == 3 {
		country, err := GetCountryByAlphaCode(webhook.Country)
		if err != nil {
			return map[string]string{}, err
		}
		webhook.Country = fmt.Sprint(country)
	}

	id := autoId()
	webhook.ID = id
	//Adds webhook to database, return documentID which will be used as webhookId
	err = db.AddToFirestore(Hash([]byte(COLLECTION)), Hash([]byte(id)), webhook)
	if err != nil {
		return map[string]string{}, err
	}
	webhooks = append(webhooks, webhook)

	//Respond with ID
	var response = make(map[string]string, 1)
	response["id"] = id

	return response, nil
}

// RunWebhookRoutine */
func RunWebhookRoutine(country string) error {
	for i, webhook := range webhooks {
		if webhook.Country == country {

			webhook.ActualCalls = webhook.ActualCalls + 1
			//Updates webhook in memory
			webhooks[i].ActualCalls = webhook.ActualCalls

			//Updates webhook in db
			err := db.UpdateWebhook(Hash([]byte(COLLECTION)), Hash([]byte(webhook.ID)), "actual_calls", webhook.ActualCalls)
			if err != nil {
				return err
			}

			if webhook.ActualCalls == webhook.Calls {
				webhook.ActualCalls = 0

				//Updates webhook in db
				err = db.UpdateWebhook(Hash([]byte(COLLECTION)), Hash([]byte(webhook.ID)), "actual_calls", webhook.ActualCalls)
				if err != nil {
					return err
				}

				//Updates webhook in memory
				webhooks[i].ActualCalls = webhook.ActualCalls

				go callUrl(webhook.Url, webhook)
			}

		}
	}
	return nil
}

// autoId Randomly generated a 15 character long string
func autoId() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	arr := make([]rune, 15)
	for i := range arr {
		arr[i] = letters[rand.Intn(len(letters))]
	}

	return string(arr)
}

// callUrl
func callUrl(url string, data interface{}) {
	payloadBuffer := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuffer).Encode(data)
	if err == nil {
		_, err = http.Post(url, "application/json", payloadBuffer)
		if err != nil {
			return
		}
	}
}
