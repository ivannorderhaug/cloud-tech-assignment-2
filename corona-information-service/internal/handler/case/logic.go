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

	if len(country) != 2 {
		country = strings.Title(strings.ToLower(country))
	}

	return country, nil
}

//Purges cache every 8 hours as the external cases API is updated three times a day
func runPurgeRoutine() {
	t = true

	ticker := time.NewTicker(12 * time.Hour)

	go func() {
		for {
			select {
			case <-ticker.C:
				cache.PurgeByDate(cases, fmt.Sprintf(time.Now().Format("2006-01-02")))
			}
		}
	}()
}

//issueGraphqlRequest uses country name to convert it into a graphql request and issues the request. Returns response
func issueGraphqlRequest(country string) (*http.Response, error) {
	query, err := graphql.ConvertToGraphql(model.QUERY, country)
	if err != nil {
		return nil, errors.New("error during marshalling")
	}

	res, err := customhttp.IssueRequest(http.MethodPost, model.CASES_URL, query)
	if err != nil {
		return nil, errors.New("error issuing the request")
	}
	return res, nil
}

//mapResponseToStruct maps the response from the graphql request into a temp struct, proceeds to map from the temp struct into the correct one.
func mapResponseToStruct(res *http.Response) (model.Case, error, int) {
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

	err := customjson.Decode(res, &tmpCase)
	if err != nil {
		return model.Case{}, errors.New("error during decoding"), http.StatusInternalServerError
	}

	if len(tmpCase.Data.Country.Name) == 0 {
		return model.Case{}, errors.New("could not find a country with that name"), http.StatusNotFound
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

	cache.Put(cases, c.Country, c)

	return c, nil, 0

}
