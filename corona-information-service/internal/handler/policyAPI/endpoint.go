package policyAPI

import (
	"corona-information-service/pkg/api"
	"corona-information-service/pkg/cache"
	"corona-information-service/tools"
	"corona-information-service/tools/customjson"
	"corona-information-service/tools/webhook"
	"fmt"
	"net/http"
)

var policies = cache.NewNestedMap()

// PolicyHandler */
func PolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}

	//Validates search
	cc, err, status := getCountryCode(r)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	//Gets date
	date, yes := hasScope(r)
	if yes {
		//Validates if date input is correctly formatted.
		if !tools.IsValidDate(date) {
			http.Error(w, "Date parameter is wrongly formatted, please see if it matches the correct format. (YYYY-MM-dd)", http.StatusBadRequest)
			return
		}
	}

	//Checks if policy with given date and alpha3 exists in cache, If it exists, it gets encoded
	if p := cache.GetNestedMap(policies, cc, date); p != nil {
		customjson.Encode(w, p)
		return
	}

	//Issues request, decodes it and returns a struct
	res, err := issueRequest(cc, date)
	if err != nil {
		http.Error(w, "Error while issuing a request", http.StatusInternalServerError)
		return
	}

	//Map data received from external api into a struct
	policy, err := mapDataToStruct(res)
	if err != nil {
		http.Error(w, "error decoding response", http.StatusInternalServerError)
		return
	}

	//Add missing data
	policy.CountryCode = cc
	policy.Scope = date

	//Adds search to cache
	cache.PutNestedMap(policies, cc, date, policy)

	//Failed webhook routine doesn't need error handling
	go func() {
		country, err := api.GetCountryNameByAlphaCode(cc)
		if err != nil {
			fmt.Println("Couldn't retrieve country name")
			return
		}
		_ = webhook.RunWebhookRoutine(fmt.Sprint(country))
	}()

	//Encodes struct
	customjson.Encode(w, policy)
}
