package handler

import (
	"corona-information-service/model"
	"corona-information-service/tools"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	path, ok, msg := tools.PathSplitter(r.URL.Path, 1)
	if !ok {
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	//Country name or isocode.
	//Converts string to lowercase before making the first letter uppercase to satisfy the graphql api search parameter
	s := strings.Title(strings.ToLower(path[0]))

	if len(s) == 3 {
		//Issues a RESTCountries api request if input is alpha3.
		//Returns the country name
		country, _ := tools.GetCountryByAlphaCode(s)
		s = fmt.Sprint(country)
	}
	query := fmt.Sprintf("query {\n  country(name: \"%s\") {\n    name\n    mostRecent {\n      date(format: \"yyyy-MM-dd\")\n      confirmed\n      recovered\n      deaths\n      growthRate\n    }\n  }\n}", s)

	jsonQuery, err := json.Marshal(model.GraphQLRequest{Query: query})
	if err != nil {
		log.Fatal(err)
	}

	res, _ := tools.IssueRequest(http.MethodPost, model.CASES_URL, jsonQuery)
	mp := unmarshalResponse(res)

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

	tools.Encode(w, c)
}

// unmarshalResponse Method for unmarshalling GraphQL response into a struct */
func unmarshalResponse(res *http.Response) model.Response {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response model.Response

	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal(err)
	}

	return response
}
