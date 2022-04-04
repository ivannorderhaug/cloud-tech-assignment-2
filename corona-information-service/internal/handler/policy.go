package handler

import (
	"corona-information-service/internal/model"
	"corona-information-service/pkg/api"
	"corona-information-service/pkg/cache"
	"corona-information-service/tools"
	"corona-information-service/tools/customhttp"
	"corona-information-service/tools/customjson"
	"corona-information-service/tools/webhook"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var policies = cache.NewNestedMap()

// PolicyHandler */
func PolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}

	path, ok := tools.PathSplitter(r.URL.Path, 1)
	if !ok {
		http.Error(w, "Path does not match the required path format specified on the root level and in the README.", http.StatusNotFound)
		return
	}

	//Alpha3 code.
	cc := strings.ToUpper(path[0])
	if len(cc) != 3 {
		http.Error(w, "Invalid alpha-3 country code. Please try again. ", http.StatusBadRequest)
		return
	}

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
	covidPolicy, err := getCovidPolicy(cc, date)
	if err != nil {
		http.Error(w, "Error while issuing a request", http.StatusInternalServerError)
		return
	}

	//Adds search to cache
	cache.PutNestedMap(policies, cc, date, covidPolicy)

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
	customjson.Encode(w, covidPolicy)
}

// getCovidPolicy Issues request to external API, decodes response into a struct, maps it correctly and returns it
func getCovidPolicy(alpha3 string, date string) (model.Policy, error) {
	url := fmt.Sprintf("%s%s/%s", model.STRINGENCY_URL, alpha3, date)

	res, err := customhttp.IssueRequest(http.MethodGet, url, nil) //returns response
	if err != nil {
		return model.Policy{}, err
	}

	var data model.TmpPolicy

	err = customjson.Decode(res, &data)
	if err != nil {
		return model.Policy{}, err
	} //returns decoded wrapper for stringency and policy data

	stringency := data.StringencyData.Stringency

	if data.StringencyData.StringencyActual != 0 {
		stringency = data.StringencyData.StringencyActual
	}

	//If there is no stringency data, the value will be set to 0 by default.
	//This changes that to -1 as to satisfy the requirements
	if stringency == 0 {
		stringency = -1
	}

	p := 0
	if len(data.PolicyActions) > 1 {
		p = len(data.PolicyActions)
	}

	return model.Policy{
		CountryCode: alpha3,
		Scope:       date,
		Stringency:  stringency,
		Policies:    p,
	}, nil
}

// hasScope Checks if date param in query exists, if not then use today's date.
func hasScope(r *http.Request) (string, bool) {
	scope := r.URL.Query().Get("scope")
	if len(scope) == 0 {
		return time.Now().Format("2006-01-02"), false
	}
	return scope, true
}
