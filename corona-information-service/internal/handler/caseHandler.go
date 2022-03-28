package handler

import (
	"corona-information-service/internal/model"
	"corona-information-service/tools"
	"fmt"
	"net/http"
	"strings"
)

// CaseHandler */
func CaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}

	path, ok := tools.PathSplitter(r.URL.Path, 1)
	if !ok {
		http.Error(w, "Path does not match the required path format specified on the root level and in the README.", http.StatusNotFound)
		return
	}

	country, err := isCountryCode(path)
	if err != nil {
		http.Error(w, "Error getting country", http.StatusInternalServerError)
		return
	}

	query, err := tools.ConvertToGraphql(model.QUERY, country)
	if err != nil {
		http.Error(w, "Error during marshalling", http.StatusInternalServerError)
		return
	}

	res, err := tools.IssueRequest(http.MethodPost, model.CASES_URL, query)
	if err != nil {
		http.Error(w, "Error issuing the request", http.StatusInternalServerError)
		return
	}

	var tmpCase model.TmpCase

	err = tools.Decode(res, &tmpCase)
	if err != nil {
		http.Error(w, "Error during decoding", http.StatusInternalServerError)
		return
	}

	if len(tmpCase.Data.Country.Name) == 0 {
		http.Error(w, "Could not find a country with that name", http.StatusNotFound)
		return
	}

	//Failed webhook routine doesn't need error handling
	go func() {
		_ = tools.RunWebhookRoutine(tmpCase.Data.Country.Name)
	}()

	info := tmpCase.Data.Country.MostRecent
	c := model.Case{
		Country:        tmpCase.Data.Country.Name,
		Date:           info.Date,
		ConfirmedCases: info.Confirmed,
		Recovered:      info.Recovered,
		Deaths:         info.Deaths,
		GrowthRate:     info.GrowthRate,
	}

	tools.Encode(w, c)
}

func isCountryCode(search []string) (string, error) {
	//Country name or alpha3.
	//Converts string to lowercase before making the first letter uppercase to satisfy the graphql api search parameter
	s := strings.Title(strings.ToLower(search[0]))
	if len(s) == 3 {
		//Issues a RESTCountries api request if input is alpha3.
		//Returns the country name
		country, err := tools.GetCountryByAlphaCode(s)
		if err != nil {
			return "", err
		}
		s = fmt.Sprint(country)
	}
	return s, nil
}
