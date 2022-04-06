package _case

import (
	"corona-information-service/internal/model"
	"corona-information-service/pkg/api"
	"corona-information-service/pkg/cache"
	"corona-information-service/tools"
	"corona-information-service/tools/customhttp"
	"corona-information-service/tools/customjson"
	"corona-information-service/tools/graphql"
	"corona-information-service/tools/webhook"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//Init empty cache
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

	country, err := getCountry(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if c := cache.Get(cases, country); c != nil {
		//Failed webhook routine doesn't need error handling
		go func() {
			_ = webhook.RunWebhookRoutine(country)
		}()
		customjson.Encode(w, c)
		return
	}

	c, err := getCase(model.CASES_URL, country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if c == nil {
		http.Error(w, "could not find a country with that name", http.StatusNotFound)
		return
	}

	//Failed webhook routine doesn't need error handling
	go func() {
		_ = webhook.RunWebhookRoutine(c.Country)
	}()

	cache.Put(cases, c.Country, c)

	customjson.Encode(w, c)
}

// getCountry handles search, converts alpha3 to country name if necessary and returns country name
func getCountry(r *http.Request) (string, error) {
	path, ok := tools.PathSplitter(r.URL.Path, 1)
	if !ok {
		return "", errors.New("path does not match the required path format specified on the root level and in the README")
	}

	//Handle spaces
	country := strings.Replace(path[0], " ", "%20", -1)

	//Gets country name if user input is alpha3 code
	if len(country) == 3 {
		alpha3ToCountry, err := api.GetCountryNameByAlphaCode(country)
		if err != nil {
			return "", errors.New("error retrieving country by country code")
		}
		country = fmt.Sprint(alpha3ToCountry)
	}

	//Handle US edge case
	if len(country) != 2 {
		country = strings.Title(strings.ToLower(country))
	}

	return country, nil
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

// getCase uses country name to issue a request, map the response into the required struct and return its reference
func getCase(url, country string) (*model.Case, error) {
	query, err := graphql.ConvertToGraphql(model.QUERY, country)
	if err != nil {
		return nil, errors.New("error during marshalling")
	}

	res, err := customhttp.IssueRequest(http.MethodPost, url, query)
	if err != nil {
		return nil, errors.New("error issuing the request")
	}

	// TmpCase Used to unwrap nested structure
	var tmpCase struct {
		Data struct {
			Country struct {
				Name       string `json:"name"`
				MostRecent struct {
					Date       string  `json:"date"`
					Confirmed  int     `json:"confirmed"`
					Recovered  int     `json:"recovered"`
					Deaths     int     `json:"deaths"`
					GrowthRate float64 `json:"growthRate"`
				} `json:"mostRecent"`
			} `json:"country"`
		} `json:"data"`
	}

	err = customjson.Decode(res, &tmpCase)
	if err != nil {
		return nil, errors.New("error during decoding")
	}

	if len(tmpCase.Data.Country.Name) == 0 {
		return nil, nil
	}

	info := tmpCase.Data.Country.MostRecent
	c := model.Case{
		Country:        tmpCase.Data.Country.Name,
		Date:           info.Date,
		ConfirmedCases: info.Confirmed,
		Recovered:      info.Recovered,
		Deaths:         info.Deaths,
		GrowthRate:     info.GrowthRate,
	}

	return &c, nil
}
