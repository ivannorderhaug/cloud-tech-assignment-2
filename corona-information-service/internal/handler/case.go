package handler

import (
	"corona-information-service/internal/model"
	"corona-information-service/tools"
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

	//Handle spaces
	country := strings.Replace(path[0], " ", "%20", -1)

	if len(country) != 2 {
		country = strings.Title(strings.ToLower(country))
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
