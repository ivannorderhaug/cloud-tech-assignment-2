package tools

import (
	"bytes"
	"corona-information-service/internal/model"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const RESTCOUNTRIES = "https://restcountries.com/v3.1/alpha/%s?fields=name"

//PathSplitter
//Example of usage: /corona/v1/cases/norway has a length of 5, if it matches basePathLength(which is 4)+length param(which is 1 in this case).
//It'll return a slice containing only 1 element which is the search param (norway)*/
func PathSplitter(path string, length int) ([]string, bool) {
	//Trims away "/" at the end of path. Only if there is one tjere
	path = strings.TrimSuffix(path, "/")
	//Splits the path into a slice, separating each part by "/"
	parts := strings.Split(path, "/")
	//Gets the length of the basePath. Length will be 4.
	basePathLength := len(strings.Split("/corona/v1/", "/"))

	if len(parts) == basePathLength {
		//Returns empty slice with an error message as the path didn't match the required format
		return []string{}, false
	}

	//Compares length of parts slice with basePath length+length param
	if len(parts) != basePathLength+length {
		//Returns empty slice with an error message as the path didn't match the required format
		return []string{}, false
	}
	return parts[basePathLength : basePathLength+length], true
}

// GetCountryByAlphaCode
// Issues a http request of method GET to the RESTCountries API
// Decodes the response and returns an interface
func GetCountryByAlphaCode(alpha3 string) (interface{}, error) {
	url := fmt.Sprintf(RESTCOUNTRIES, alpha3)
	// Create new request
	res, err := IssueRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	//Local struct, only used for this purpose
	var c struct {
		Name interface{} `json:"name"`
	}

	err = Decode(res, &c)
	if err != nil {
		return nil, err
	}

	//Returns an interface by going one layer into the country name interface and picking out the common name
	return c.Name.(map[string]interface{})["common"], nil
}

//IsValidDate Uses Regular Expressions to validate if string matches required format */
func IsValidDate(date string) bool {
	//YYYY-mm-dd
	pattern := regexp.MustCompile("([12]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01]))")
	return pattern.MatchString(date)
}

// IssueRequest */
func IssueRequest(method string, url string, body []byte) (*http.Response, error) {
	// Create new request
	r, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return &http.Response{}, err
	}
	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")

	// Instantiate the client
	client := &http.Client{}

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		return &http.Response{}, err
	}

	return res, nil
}

// RemoveIndex */
func RemoveIndex(s []model.Webhook, index int) []model.Webhook {
	ret := make([]model.Webhook, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}
