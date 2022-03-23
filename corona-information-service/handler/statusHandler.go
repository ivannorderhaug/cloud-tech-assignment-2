package handler

import (
	"corona-information-service/model"
	"corona-information-service/tools"
	"encoding/json"
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

	//This has to be done as the casesAPI only accepts POST
	query, _ := json.Marshal(model.GraphQLRequest{Query: model.QUERY})
	casesApi, err := tools.IssueRequest(http.MethodPost, model.CASES_URL, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	policyApi, err := tools.IssueRequest(http.MethodHead, model.STRINGENCY_URL+"nor/"+"2022-02-04", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	restCountriesApi, err := tools.IssueRequest(http.MethodHead, model.RESTCOUNTRIES_URL, nil)
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
