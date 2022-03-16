package handler

import (
	"corona-information-service/functions"
	"corona-information-service/model"
	"encoding/json"
	"fmt"
	"net/http"
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

	country := path[0]

	query := fmt.Sprintf("query {\n  country(name: \"%s\") {\n    name\n    mostRecent {\n      date(format: \"yyyy-MM-dd\")\n      confirmed\n      recovered\n      deaths\n      growthRate\n    }\n  }\n}", country)

	jsonQuery, _ := json.Marshal(model.GraphQLRequest{Query: query})

	_ = functions.IssueGraphQLRequest(url, jsonQuery)
}
