package webhook

import (
	"bytes"
	"corona-information-service/internal/model"
	"corona-information-service/pkg/api"
	"corona-information-service/pkg/customjson"
	"corona-information-service/pkg/db"
	"corona-information-service/tools/hash"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const COLLECTION = "notifications"

var webhooks []model.Webhook

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
	documentsFromFirestore, err := db.GetAllDocumentsFromFirestore(hash.Hash(COLLECTION))
	if err != nil {
		return []model.Webhook{}, err
	}
	webhooks = make([]model.Webhook, 0)
	//Converts each document snapshot into a webhook interface and adds it to the global webhooks slice
	for _, documentSnapshot := range documentsFromFirestore {
		var webhook model.Webhook
		err = documentSnapshot.DataTo(&webhook)
		if err != nil {
			return []model.Webhook{}, err
		}
		webhooks = append(webhooks, webhook)
	}
	return webhooks, nil
}

// DeleteWebhook */
func DeleteWebhook(webhookId string) (bool, error) {
	var deleted = false
	if len(webhooks) != 0 {
		for i, wh := range webhooks {
			if wh.ID == webhookId {
				copy(webhooks[i:], webhooks[i+1:])          // Shift s[i+1:] left one index.
				webhooks[len(webhooks)-1] = model.Webhook{} // Erase last element.
				webhooks = webhooks[:len(webhooks)-1]
				deleted = true
			}
		}
	}

	if err := db.DeleteSingleDocumentFromFirestore(hash.Hash(COLLECTION), hash.Hash(webhookId)); err != nil {
		return false, err
	}

	return deleted, nil
}

// RegisterWebhook */
func RegisterWebhook(r *http.Request) (map[string]string, error) {
	var webhook model.Webhook

	err := customjson.Decode(r, &webhook)
	if err != nil {
		return map[string]string{}, err
	}

	//checks if alpha3 code was used as param for country
	if len(webhook.Country) == 3 {
		country, err := api.GetCountryNameByAlphaCode(webhook.Country)
		if err != nil {
			return map[string]string{}, err
		}
		webhook.Country = fmt.Sprint(country)
	}

	id := autoId()
	webhook.ID = id
	//Adds webhook to database, return documentID which will be used as webhookId
	err = db.AddToFirestore(hash.Hash(COLLECTION), hash.Hash(id), webhook)
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
			err := db.UpdateDocument(hash.Hash(COLLECTION), hash.Hash(webhook.ID), "actual_calls", webhook.ActualCalls)
			if err != nil {
				return err
			}

			if webhook.ActualCalls == webhook.Calls {
				webhook.ActualCalls = 0

				//Updates webhook in db
				err = db.UpdateDocument(hash.Hash(COLLECTION), hash.Hash(webhook.ID), "actual_calls", webhook.ActualCalls)
				if err != nil {
					return err
				}

				//Updates webhook in memory
				webhooks[i].ActualCalls = webhook.ActualCalls

				webhook.Invoked = time.Now().UTC().String()

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
