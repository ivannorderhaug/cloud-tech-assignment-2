package handler

import (
	"corona-information-service/internal/model"
	"corona-information-service/pkg/customhttp"
	"corona-information-service/pkg/customjson"
	"corona-information-service/pkg/webhook"
	"fmt"
	"net/http"
	"time"
)

var startTime = time.Now()

//getUptime: Gets uptime
func getUptime() time.Duration {
	return time.Since(startTime)
}

// StatusHandler */
func StatusHandler(client customhttp.HTTPClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
			return
		}
		//Declare vars
		var casesApiStatus, policyApiStatus, restCountriesApiStatus string

		//Requests
		casesApi, err := customhttp.IssueRequest(client, http.MethodGet, model.CASES_API, nil)
		if err != nil {
			casesApiStatus = fmt.Sprintf("%s %s", http.StatusFailedDependency, http.StatusText(http.StatusFailedDependency))
		}
		policyApi, err := customhttp.IssueRequest(client, http.MethodHead, model.STRINGENCY_API, nil)
		if err != nil {
			policyApiStatus = fmt.Sprintf("%s %s", http.StatusFailedDependency, http.StatusText(http.StatusFailedDependency))
		}
		restCountriesApi, err := customhttp.IssueRequest(client, http.MethodGet, model.RESTCOUNTRIES_API, nil)
		if err != nil {
			restCountriesApiStatus = fmt.Sprintf("%s %s", http.StatusFailedDependency, http.StatusText(http.StatusFailedDependency))
		}

		//Statuses
		casesApiStatus = casesApi.Status
		policyApiStatus = policyApi.Status
		restCountriesApiStatus = restCountriesApi.Status

		webhooksCount := 0
		webhooks, err := webhook.GetAllWebhooks()
		if err == nil {
			webhooksCount = len(webhooks)
		}

		status := model.Status{
			CasesApi:      casesApiStatus,
			PolicyApi:     policyApiStatus,
			RestCountries: restCountriesApiStatus,
			Webhooks:      webhooksCount,
			Version:       model.VERSION,
			Uptime:        int(getUptime().Seconds()),
		}

		customjson.Encode(w, &status)

	}
}
