package handler

import (
	"corona-information-service/internal/model"
	"corona-information-service/tools"
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

	casesApi, err := tools.IssueRequest(http.MethodGet, model.CASES_API, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	policyApi, err := tools.IssueRequest(http.MethodHead, model.STRINGENCY_API, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	restCountriesApi, err := tools.IssueRequest(http.MethodHead, model.RESTCOUNTRIES_API, nil)
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

	tools.Encode(w, status)

}
