package handler

import (
	"corona-information-service/model"
	"corona-information-service/tools"
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
		http.Error(w, "Error while issuing a request", http.StatusInternalServerError)
		return
	}

	//Encodes struct
	tools.Encode(w, covidPolicy)
}

func getCovidPolicy(alpha3 string, date string) (model.Policy, error) {
	url := tools.MakeURL(model.STRINGENCY_URL, alpha3, date)

	res, err := tools.IssueRequest(http.MethodGet, url, nil) //returns response
	if err != nil {
		return model.Policy{}, err
	}

	var data model.TmpPolicy

	err2 := tools.Decode(res, &data)
	if err2 != nil {
		return model.Policy{}, err2
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
