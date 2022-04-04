package handler

import (
	"corona-information-service/internal/model"
	"corona-information-service/tools/customhttp"
	"corona-information-service/tools/customjson"
	"corona-information-service/tools/webhook"
	"net/http"
	"time"
)

var startTime = time.Now()

//getUptime: Gets uptime
func getUptime() time.Duration {
	return time.Since(startTime)
}

// StatusHandler */
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}

	casesApi, err := customhttp.IssueRequest(http.MethodGet, model.CASES_API, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	policyApi, err := customhttp.IssueRequest(http.MethodHead, model.STRINGENCY_API, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	restCountriesApi, err := customhttp.IssueRequest(http.MethodHead, model.RESTCOUNTRIES_API, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	webhooksCount := 0
	webhooks, err := webhook.GetAllWebhooks()
	if err == nil {
		webhooksCount = len(webhooks)
	}

	status := model.Status{
		CasesApi:      casesApi.Status,
		PolicyApi:     policyApi.Status,
		RestCountries: restCountriesApi.Status,
		Webhooks:      webhooksCount,
		Version:       model.VERSION,
		Uptime:        int(getUptime().Seconds()),
	}

	customjson.Encode(w, status)

}
