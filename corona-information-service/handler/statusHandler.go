package handler

import (
	"corona-information-service/model"
	"encoding/json"
	"net/http"
	"time"
)

var startTime = time.Now()

//getUptime: Gets getUptime
func getUptime() time.Duration {
	return time.Since(startTime)
}

// StatusHandler */
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}

	casesApi, err := getStatus(CASES_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	policyApi, err := getStatus(STRINGENCY_URL + "nor/" + "2022-02-04")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	restCountriesApi, err := getStatus(RESTCOUNTRIES_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	status := model.Status{
		CasesApi:      casesApi,
		PolicyApi:     policyApi,
		RestCountries: restCountriesApi,
		Version:       VERSION,
		Uptime:        int(getUptime().Seconds()),
	}

	encodeStatusInformation(w, status)

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

// encodeStatusInformation */
func encodeStatusInformation(w http.ResponseWriter, r model.Status) {
	// Write content type header
	w.Header().Add("content-type", "application/json")

	// Instantiate encoder
	encoder := json.NewEncoder(w)

	//Encodes response
	err := encoder.Encode(r)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}
}
