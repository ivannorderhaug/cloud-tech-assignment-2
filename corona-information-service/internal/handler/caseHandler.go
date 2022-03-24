package handler

import (
	model2 "corona-information-service/internal/model"
	tools2 "corona-information-service/internal/tools"
	"encoding/json"
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

	path, ok, msg := tools2.PathSplitter(r.URL.Path, 1)
	if !ok {
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	//Country name or alpha3.
	//Converts string to lowercase before making the first letter uppercase to satisfy the graphql api search parameter
	s := strings.Title(strings.ToLower(path[0]))

	if len(s) == 2 {
		s = strings.ToUpper(s)
	}

	if len(s) == 3 {
		//Issues a RESTCountries api request if input is alpha3.
		//Returns the country name
		country, _ := tools2.GetCountryByAlphaCode(s)
		s = fmt.Sprint(country)

	}

	//Formats query
	query := fmt.Sprintf(model2.QUERY, s)

	fmt.Println(query)
	jsonQuery, err := json.Marshal(model2.GraphQLRequest{Query: query})
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}

	res, _ := tools2.IssueRequest(http.MethodPost, model2.CASES_URL, jsonQuery)

	var tmpCase model2.TmpCase
	decode := tools2.Decode(res, &tmpCase)
	if decode != nil {
		http.Error(w, "Error during decoding", http.StatusInternalServerError)
		return
	}

	if len(tmpCase.Data.Country.Name) == 0 {
		http.Error(w, "Could not find a country with that name", http.StatusNotFound)
		return
	}

	info := tmpCase.Data.Country.MostRecent
	c := model2.Case{
		Country:        tmpCase.Data.Country.Name,
		Date:           info.Date,
		ConfirmedCases: info.Confirmed,
		Recovered:      info.Recovered,
		Deaths:         info.Deaths,
		GrowthRate:     info.GrowthRate,
	}

	tools2.Encode(w, c)
}
