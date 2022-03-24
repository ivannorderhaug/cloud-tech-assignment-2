package handler

import (
	"corona-information-service/internal/model"
	tools2 "corona-information-service/tools"
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

	casesApi, err := tools2.IssueRequest(http.MethodGet, model.CASES_URL+"?query=%7B__typename%7D", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	policyApi, err := tools2.IssueRequest(http.MethodHead, model.STRINGENCY_URL+"nor/"+"2022-02-04", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	restCountriesApi, err := tools2.IssueRequest(http.MethodHead, model.RESTCOUNTRIES_URL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	status := model.Status{
		CasesApi:      casesApi.Status,
		PolicyApi:     policyApi.Status,
		RestCountries: restCountriesApi.Status,
		Version:       model.VERSION,
		Uptime:        int(getUptime().Seconds()),
	}

	tools2.Encode(w, status)

}
