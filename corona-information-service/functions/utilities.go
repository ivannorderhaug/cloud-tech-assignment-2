package functions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

//PathSplitter */
func PathSplitter(path string, length int) ([]string, bool, string) {
	//Trims away last "/"
	path = strings.TrimSuffix(path, "/")
	//Splits the path into a slice, separating each part by "/"
	parts := strings.Split(path, "/")
	//Gets the length of the basePath
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
func GetCountryByAlphaCode(isocode string) interface{} {
	url := fmt.Sprintf("https://restcountries.com/v3.1/alpha/%s?fields=name", isocode)

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
