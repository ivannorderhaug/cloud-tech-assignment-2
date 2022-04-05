package _policy

import (
	"corona-information-service/internal/model"
	"corona-information-service/pkg/api"
	"corona-information-service/pkg/cache"
	"corona-information-service/tools"
	"corona-information-service/tools/customhttp"
	"corona-information-service/tools/customjson"
	"corona-information-service/tools/webhook"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var policies = cache.NewNestedMap()

// issueRequest Issues request to external API, decodes response into a struct, maps it correctly and returns it
func issueRequest(alpha3 string, date string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s/%s", model.STRINGENCY_URL, alpha3, date)

	//Issues request, gets response
	res, err := customhttp.IssueRequest(http.MethodGet, url, nil) //returns response
	if err != nil {
		return nil, err
	}
	return res, nil
}

//mapDataToStruct maps the received data into the struct
func mapDataToStruct(res *http.Response) (model.Policy, error) {
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
		return model.Policy{}, err
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
	p := 0
	if len(data.PolicyActions) > 1 {
		p = len(data.PolicyActions)
	}

	return model.Policy{
		Stringency: stringency,
		Policies:   p,
	}, nil
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
