package cases

import (
	"corona-information-service/internal/model"
	"corona-information-service/pkg/api"
	"corona-information-service/pkg/cache"
	"corona-information-service/pkg/customhttp"
	"corona-information-service/pkg/customjson"
	"corona-information-service/pkg/webhook"
	"corona-information-service/tools"
	"corona-information-service/tools/graphql"
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
func CaseHandler(client customhttp.HTTPClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
			return
		}

		//Runs purge routine in background
		if !t {
			runPurgeRoutine()
		}

		//Gets correct country name from the request issued by the end user
		country, err, status := getCountry(r)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		//Checks cache if case data for country already exists
		if c := cache.Get(cases, country); c != nil {
			//Failed webhook routine doesn't need error handling
			go func() {
				_ = webhook.RunWebhookRoutine(country)
			}()
			customjson.Encode(w, c)
			return
		}

		//Formats the country correctly into a graphql query
		query, err := graphql.ConvertToGraphql(model.QUERY, country)
		if err != nil {
			http.Error(w, "error during marshalling", http.StatusInternalServerError)
			return
		}

		//Issues the graphql request
		res, err := customhttp.IssueRequest(client, http.MethodPost, model.CASES_URL, query)
		if err != nil {
			http.Error(w, "error issuing a request", http.StatusInternalServerError)
			return
		}

		//Gets case data from response
		c, err, status := getCase(res)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		//Runs webhook routine on a different thread
		//Failed webhook routine doesn't need error handling
		go func() {
			_ = webhook.RunWebhookRoutine(c.Country)
		}()

		//At this point, it is clear that the case data does not exist in the cache, therefore add it to the cache
		cache.Put(cases, c.Country, c)

		//Encode case data
		customjson.Encode(w, c)
	}
}

// getCountry handles search, converts alpha3 to country name if necessary and returns country name
func getCountry(r *http.Request) (string, error, int) {
	//Checks if path is correctly formatted
	path, ok := tools.PathSplitter(r.URL.Path, 1)
	if !ok {
		return "", errors.New("path does not match the required path format specified on the root level and in the README"), http.StatusBadRequest
	}

	//Count spaces, due to some edge cases in external api
	count := 0
	for _, v := range path[0] {
		if v == ' ' {
			count++
		}
	}

	//Handle spaces
	country := path[0]

	//Gets country name if user input is alpha3 code
	if len(country) == 3 {
		alpha3ToCountry, err := api.GetCountryNameByAlphaCode(country)
		if err != nil {
			return "", errors.New("error retrieving country by country code"), http.StatusInternalServerError
		}
		country = fmt.Sprint(alpha3ToCountry)
	}

	//If space count is lower than 2, then it is safe to assume that the first letter should be uppercase,
	//if there are multiple words then the first letter will be uppercase in each word
	if count < 2 {
		country = strings.Title(strings.ToLower(country))
	}

	//Handle US edge case
	if len(country) == 2 {
		country = strings.ToUpper(country)
	}

	return country, nil, 0
}

//Purges cache every 8 hours as the external case API is updated three times a day
func runPurgeRoutine() {
	t = true
	//Create new ticker
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
func getCase(res *http.Response) (*model.Case, error, int) {
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

	//Decodes to struct
	err := customjson.Decode(res, &tmpCase)
	if err != nil {
		return &model.Case{}, errors.New("error during decoding"), http.StatusInternalServerError
	}

	//If name in struct is empty, then the request failed
	if len(tmpCase.Data.Country.Name) == 0 {
		return &model.Case{}, errors.New("could not find a country with that name"), http.StatusNotFound
	}

	//Map data
	info := tmpCase.Data.Country.MostRecent
	c := model.Case{
		Country:        tmpCase.Data.Country.Name,
		Date:           info.Date,
		ConfirmedCases: info.Confirmed,
		Recovered:      info.Recovered,
		Deaths:         info.Deaths,
		GrowthRate:     info.GrowthRate,
	}

	return &c, nil, 0
}
