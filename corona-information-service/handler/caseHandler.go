package handler

import (
	"corona-information-service/functions"
	"corona-information-service/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// CaseHandler */
func CaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}

	// URL to invoke
	url := "https://covid19-graphql.vercel.app/"

	path, ok, msg := functions.PathSplitter(r.URL.Path, 1)
	if !ok {
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	//Country name or isocode.
	//Converts string to lowercase before making the first letter uppercase to satisfy the graphql api search parameter
	s := strings.Title(strings.ToLower(path[0]))

	if len(s) == 3 {
		//Issues a RESTCountries api request if input is isocode.
		//Returns the country name
		s = fmt.Sprint(functions.GetCountryByAlphaCode(s))
	}
	query := fmt.Sprintf("query {\n  country(name: \"%s\") {\n    name\n    mostRecent {\n      date(format: \"yyyy-MM-dd\")\n      confirmed\n      recovered\n      deaths\n      growthRate\n    }\n  }\n}", s)

	jsonQuery, err := json.Marshal(model.GraphQLRequest{Query: query})
	if err != nil {
		log.Fatal(err)
	}

	res := functions.IssueGraphQLRequest(url, jsonQuery)
	mp := functions.UnmarshalResponse(res)

	if len(mp.Data.Country.Name) == 0 {
		http.Error(w, "Could not find a country with that name", http.StatusNotFound)
		return
	}

	c := model.Case{
		Country:        mp.Data.Country.Name,
		Date:           mp.Data.Country.Info.Date,
		ConfirmedCases: mp.Data.Country.Info.Confirmed,
		Recovered:      mp.Data.Country.Info.Recovered,
		Deaths:         mp.Data.Country.Info.Deaths,
		GrowthRate:     mp.Data.Country.Info.GrowthRate,
	}

	functions.EncodeCaseInformation(w, c)

}
