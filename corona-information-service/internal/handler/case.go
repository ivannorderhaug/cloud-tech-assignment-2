package handler

import (
	"corona-information-service/internal/model"
	"corona-information-service/pkg/api"
	"corona-information-service/pkg/cache"
	"corona-information-service/tools"
	"corona-information-service/tools/customhttp"
	"corona-information-service/tools/customjson"
	"corona-information-service/tools/graphql"
	"corona-information-service/tools/webhook"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var cases = cache.New()

//Bool used to make sure the purge routine is only run once
var t = false

// CaseHandler */
func CaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}

	if !t {
		runPurgeRoutine()
	}

	path, ok := tools.PathSplitter(r.URL.Path, 1)
	if !ok {
		http.Error(w, "Path does not match the required path format specified on the root level and in the README.", http.StatusNotFound)
		return
	}

	//Handle spaces
	country := strings.Replace(path[0], " ", "%20", -1)

	if len(country) == 3 {
		alpha3ToCountry, err := api.GetCountryNameByAlphaCode(country)
		if err != nil {
			http.Error(w, "Error retrieving country by country code", http.StatusInternalServerError)
			return
		}
		country = fmt.Sprint(alpha3ToCountry)
	}

	if len(country) != 2 {
		country = strings.Title(strings.ToLower(country))
	}

	if c := cache.Get(cases, country); c != nil {
		customjson.Encode(w, c)
		return
	}

	query, err := graphql.ConvertToGraphql(model.QUERY, country)
	if err != nil {
		http.Error(w, "Error during marshalling", http.StatusInternalServerError)
		return
	}

	res, err := customhttp.IssueRequest(http.MethodPost, model.CASES_URL, query)
	if err != nil {
		http.Error(w, "Error issuing the request", http.StatusInternalServerError)
		return
	}

	var tmpCase model.TmpCase

	err = customjson.Decode(res, &tmpCase)
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
		_ = webhook.RunWebhookRoutine(tmpCase.Data.Country.Name)
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

	cache.Put(cases, country, c)
	customjson.Encode(w, c)
}

//Purges cache every 8 hours as the external case API is updated three times a day
func runPurgeRoutine() {
	t = true
	ticker := time.NewTicker(8 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				cache.PurgeByDate(cases, fmt.Sprintf(time.Now().Format("2006-01-02")))
			}
		}
	}()
}
