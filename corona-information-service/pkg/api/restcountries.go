package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var RESTCOUNTRIES = "https://restcountries.com/v3.1/alpha/%s?fields=name"

// GetCountryNameByAlphaCode
// Issues a http request of method GET to the RESTCountries API
// Decodes the response and returns an interface
func GetCountryNameByAlphaCode(alpha3 string) (string, error) {
	url := fmt.Sprintf(RESTCOUNTRIES, alpha3)
	// Create new request
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	//Local struct, only used for this purpose
	var c struct {
		Name interface{} `json:"name"`
	}

	decoder := json.NewDecoder(res.Body)

	err = decoder.Decode(&c)
	if err != nil {
		return "", err
	}

	if c.Name == nil {
		return "", errors.New("country code does not exist")
	}

	//Returns an interface by going one layer into the country name interface and picking out the common name
	return fmt.Sprint(c.Name.(map[string]interface{})["common"]), nil
}
