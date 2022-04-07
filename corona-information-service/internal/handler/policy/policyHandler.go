package _policy

import (
	"corona-information-service/internal/model"
	"corona-information-service/pkg/api"
	"corona-information-service/pkg/cache"
	"corona-information-service/pkg/customhttp"
	"corona-information-service/pkg/customjson"
	"corona-information-service/pkg/webhook"
	"corona-information-service/tools"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var policies = cache.NewNestedMap()

// PolicyHandler */
func PolicyHandler(client customhttp.HTTPClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
			return
		}

		//Validates search
		cc, err, status := getCountryCode(r)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		//Gets date
		date, yes := hasScope(r)
		if yes {
			//Validates if date input is correctly formatted.
			if !tools.IsValidDate(date) {
				http.Error(w, "Date parameter is wrongly formatted, please see if it matches the correct format. (YYYY-MM-dd)", http.StatusBadRequest)
				return
			}
		}

		//Checks if policy with given date and alpha3 exists in cache, If it exists, it gets encoded
		if cachedP := cache.GetNestedMap(policies, cc, date); cachedP != nil {
			runWebhookRoutine(cc)
			customjson.Encode(w, cachedP)
			return
		}

		//URL TO INVOKE
		url := fmt.Sprintf("%s%s/%s", model.STRINGENCY_URL, cc, date)

		//Issues request, gets response
		res, err := customhttp.IssueRequest(client, http.MethodGet, url, nil) //returns response
		if err != nil {
			http.Error(w, "Error while issuing a request", http.StatusInternalServerError)
			return
		}

		//Map data received from external api into a struct
		p, err := mapDataToStruct(res, cc, date)
		if err != nil {
			http.Error(w, "error decoding response", http.StatusInternalServerError)
			return
		}

		//Adds search to cache
		cache.PutNestedMap(policies, cc, date, p)

		//Failed webhook routine doesn't need error handling
		runWebhookRoutine(cc)

		//Encodes struct
		customjson.Encode(w, p)
	}
}

//mapDataToStruct maps the received data into the struct
func mapDataToStruct(res *http.Response, cc, date string) (*model.Policy, error) {
	// Used to unwrap nested structure
	var data struct {
		StringencyData struct {
			Stringency       float64 `json:"stringency"`
			StringencyActual float64 `json:"stringency_actual,omitempty"`
		} `json:"stringencyData"`
		PolicyActions []interface{} `json:"policyActions"`
	}

	err := customjson.Decode(res, &data)
	if err != nil {
		return &model.Policy{}, err
	} //returns decoded wrapper for stringency and policy data

	stringency := data.StringencyData.Stringency

	if data.StringencyData.StringencyActual != 0 {
		stringency = data.StringencyData.StringencyActual
	}

	//If there is no stringency data, the value will be set to 0 by default.
	//This changes that to -1 as to satisfy the requirements
	if stringency == 0 {
		stringency = -1
	}

	//Assumption: Active policies are the number of policies returned.
	pa := 0
	if len(data.PolicyActions) > 1 {
		pa = len(data.PolicyActions)
	}

	p := model.Policy{
		CountryCode: cc,
		Scope:       date,
		Stringency:  stringency,
		Policies:    pa,
	}

	return &p, nil
}

// hasScope Checks if date param in query exists, if not then use today's date.
func hasScope(r *http.Request) (string, bool) {
	scope := r.URL.Query().Get("scope")
	if len(scope) == 0 {
		return time.Now().Format("2006-01-02"), false
	}
	return scope, true
}

func getCountryCode(r *http.Request) (string, error, int) {
	path, ok := tools.PathSplitter(r.URL.Path, 1)
	if !ok {
		return "", errors.New("path does not match the required path format specified on the root level and in the README"), http.StatusNotFound
	}

	//Alpha3 code.
	cc := strings.ToUpper(path[0])
	if len(cc) != 3 {
		return "", errors.New("invalid alpha-3 country code. Please try again"), http.StatusBadRequest
	}

	return cc, nil, 0
}

func runWebhookRoutine(cc string) {
	go func() {
		country, err := api.GetCountryNameByAlphaCode(cc)
		if err != nil {
			fmt.Println("Couldn't retrieve country name")
			return
		}
		_ = webhook.RunWebhookRoutine(fmt.Sprint(country))
	}()
}
