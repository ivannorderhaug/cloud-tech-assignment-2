package functions

import (
	"fmt"
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
	//Compares length of parts slice with basePath length+length param
	if len(parts) != basePathLength+length {
		//Returns empty slice with an error message as the path didn't match the required format
		return []string{}, false, "Path not found, make sure the path matches the required path format specified on the root level and in the README."
	}
	return parts[basePathLength : basePathLength+length], true, "" //Empty message
}

//IssueGraphQLRequest Issues a http request of method POST. Returns response */
func IssueGraphQLRequest(url string, jsonQuery []byte) *http.Response {
	// Create new request
	r, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(jsonQuery)))
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
	defer res.Body.Close()

	return res
}
