package handler

import (
	"corona-information-service/model"
	"corona-information-service/tools"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// PolicyHandler */
func PolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}

	path, ok, msg := tools.PathSplitter(r.URL.Path, 1)
	if !ok {
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	//Country name or alpha3 code.
	s := strings.ToUpper(path[0])
	if len(s) != 3 {
		http.Error(w, "Invalid alpha-3 country code. Please try again. ", http.StatusBadRequest)
		return
	}
	//Checks if date param in query exists, if not then use todays date.
	date := r.URL.Query().Get("scope")
	if len(date) == 0 {
		date = time.Now().Format("2006-01-02")
	}

	//Validates if date input is correctly formatted.
	if !tools.IsValidDate(date) {
		http.Error(w, "Date parameter is wrongly formatted, please see if it matches the correct format. (YYYY-MM-dd)", http.StatusBadRequest)
		return
	}

	//Issues request, decodes it and returns a struct
	covidPolicy, err := getCovidPolicy(s, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	//Encodes struct
	encodePolicyInformation(w, covidPolicy)
}

func getCovidPolicy(alpha3 string, date string) (model.Policy, error) {
	url := STRINGENCY_URL + alpha3 + "/" + date

	res, err := issueRequest(url) //returns response
	if err != nil {
		return model.Policy{}, err
	}

	w, err := decodeResponse(res) //returns decoded wrapper for stringency data
	if err != nil {
		return model.Policy{}, err
	}

	stringency := w.StringencyData.Stringency

	if w.StringencyData.StringencyActual != 0 {
		stringency = w.StringencyData.StringencyActual
	}

	//If there is no stringency data, the value will be set to 0 by default.
	//This changes that to -1 as to satisfy the requirements
	if stringency == 0 {
		stringency = -1
	}

	if len(w.PolicyActions) > 1 {
		return model.Policy{
			CountryCode: alpha3,
			Scope:       date,
			Stringency:  stringency,
			Policies:    len(w.PolicyActions),
		}, nil
	} else {
		return model.Policy{
			CountryCode: alpha3,
			Scope:       date,
			Stringency:  stringency,
			Policies:    0,
		}, nil
	}
}

// issueRequest */
func issueRequest(url string) (*http.Response, error) {
	// Create new request
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return &http.Response{}, err
	}
	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")

	// Instantiate the client
	client := &http.Client{}

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		return &http.Response{}, err
	}

	return res, nil
}

// decodeResponse */
func decodeResponse(res *http.Response) (model.CovidPolicyData, error) {
	var w model.CovidPolicyData

	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&w); err != nil {
		return model.CovidPolicyData{}, err
	}

	return w, nil
}

// encodePolicyInformation */
func encodePolicyInformation(w http.ResponseWriter, r model.Policy) {
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
