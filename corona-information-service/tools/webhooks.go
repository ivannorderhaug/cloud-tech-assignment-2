package tools

import (
	"bytes"
	"corona-information-service/internal/db"
	"corona-information-service/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
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
		documentsFromFirestore, err := db.GetAllDocumentsFromFirestore(COLLECTION)
		if err != nil {
			return []model.Webhook{}, err
		}

		//Converts each document snapshot into a webhook interface and adds it to the global webhooks slice
		for _, documentSnapshot := range documentsFromFirestore {
			var webhook model.Webhook
			webhook.ID = documentSnapshot.Ref.ID
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

	if err := db.DeleteSingleDocumentFromFirestore(COLLECTION, webhookId); err != nil {
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

	webhookID, err := db.AddToFirestore(COLLECTION, webhook)
	if err != nil {
		return map[string]string{}, err
	}

	webhook.ID = webhookID
	webhooks = append(webhooks, webhook)

	//Respond with ID
	var response = make(map[string]string, 1)
	response["id"] = webhookID

	return response, nil
}

func RunWebhookRoutine(country string) error {
	for i, webhook := range webhooks {
		if webhook.Country == country {

			webhook.ActualCalls = webhook.ActualCalls + 1

			err := db.UpdateWebhook(COLLECTION, webhook.ID, webhook.ActualCalls)
			if err != nil {
				return err
			}

			//This removes the webhook from memory
			webhooks = RemoveIndex(webhooks, i)

			//This adds the webhook back into the cache, with updated data
			webhooks = append(webhooks, webhook)

			if webhook.ActualCalls == webhook.Calls {
				webhook.ActualCalls = 0
				//This removes the webhook from memory
				webhooks = RemoveIndex(webhooks, i)

				//This adds the webhook back into the cache, with updated data
				webhooks = append(webhooks, webhook)

				err = callUrl(webhook.Url, webhook)
				if err != nil {
					return err
				}
			}

		}
	}
	return nil
}

// callUrl
func callUrl(url string, data interface{}) error {
	payloadBuffer := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuffer).Encode(data)
	if err == nil {
		_, err = http.Post(url, "application/json", payloadBuffer)
		if err != nil {
			return err
		}
	}
	return nil
}
