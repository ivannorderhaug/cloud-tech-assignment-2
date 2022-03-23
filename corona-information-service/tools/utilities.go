package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

//PathSplitter
//Example of usage: /corona/v1/cases/norway has a length of 5, if it matches basePathLength(which is 4)+length param(which is 1 in this case).
//It'll return a slice containing only 1 element which is the search param (norway)*/
func PathSplitter(path string, length int) ([]string, bool, string) {
	//Trims away "/" at the end of path. Only if there is one tjere
	path = strings.TrimSuffix(path, "/")
	//Splits the path into a slice, separating each part by "/"
	parts := strings.Split(path, "/")
	//Gets the length of the basePath. Length will be 4.
	basePathLength := len(strings.Split("/corona/v1/", "/"))

	if len(parts) == basePathLength {
		//Returns empty slice with an error message as the path didn't match the required format
		return []string{}, false, "Missing search parameter"
	}

	//Compares length of parts slice with basePath length+length param
	if len(parts) != basePathLength+length {
		//Returns empty slice with an error message as the path didn't match the required format
		return []string{}, false, "Path not found, make sure the path matches the required path format specified on the root level and in the README."
	}
	return parts[basePathLength : basePathLength+length], true, "" //Empty message
}

// GetCountryByAlphaCode
// Issues a http request of method GET to the RESTCountries API
// Decodes the response and returns an interface
func GetCountryByAlphaCode(alpha3 string) interface{} {
	url := fmt.Sprintf("https://restcountries.com/v3.1/alpha/%s?fields=name", alpha3)

	// Create new request
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Errorf("Error in creating request:", err.Error())
	}
	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")

	// Instantiate the client
	client := &http.Client{}

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		fmt.Errorf("Error in response:", err.Error())
	}

	//Creates decoder to decode response from GET request
	decoder := json.NewDecoder(res.Body)

	//Local struct, only used for this purpose
	var c struct {
		Name interface{} `json:"name"`
	}

	//Populates slice if decoder is successful
	if err := decoder.Decode(&c); err != nil {
		log.Fatal(err)
	}

	//Returns an interface by going one layer into the country name interface and picking out the common name
	return c.Name.(map[string]interface{})["common"]
}

//IsValidDate Uses Regular Expressions to validate if string matches required format */
func IsValidDate(date string) bool {
	//YYYY-mm-dd
	pattern := regexp.MustCompile("([12]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01]))")
	return pattern.MatchString(date)
}
