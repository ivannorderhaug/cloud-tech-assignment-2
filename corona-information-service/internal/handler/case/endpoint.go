package _case

import (
	"corona-information-service/pkg/cache"
	"corona-information-service/tools/customjson"
	"corona-information-service/tools/webhook"
	"net/http"
)

// CaseHandler */
func CaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}

	if !t {
		runPurgeRoutine()
	}

	country, err := getCountry(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if c := cache.Get(cases, country); c != nil {
		//Failed webhook routine doesn't need error handling
		go func() {
			_ = webhook.RunWebhookRoutine(country)
		}()
		customjson.Encode(w, c)
		return
	}

	res, err := issueGraphqlRequest(country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c, err, status := mapResponseToStruct(res)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	//Failed webhook routine doesn't need error handling
	go func() {
		_ = webhook.RunWebhookRoutine(c.Country)
	}()

	cache.Put(cases, c.Country, c)

	customjson.Encode(w, c)
}
