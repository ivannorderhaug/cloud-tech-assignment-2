package handler

import (
	"corona-information-service/model"
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

	casesApi, err := getStatus(model.CASES_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	policyApi, err := getStatus(model.STRINGENCY_URL + "nor/" + "2022-02-04")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	restCountriesApi, err := getStatus(model.RESTCOUNTRIES_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	status := model.Status{
		CasesApi:      casesApi,
		PolicyApi:     policyApi,
		RestCountries: restCountriesApi,
		Version:       model.VERSION,
		Uptime:        int(getUptime().Seconds()),
	}

	tools.Encode(w, status)

}

//getStatus Simple method to retrieve a status code from an external api
func getStatus(api string) (string, error) {
	// Create new request
	r, err := http.NewRequest(http.MethodHead, api, nil)
	if err != nil {
		return "", err
	}
	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")

	// Instantiate the client
	client := &http.Client{}

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}
	return res.Status, nil
}
