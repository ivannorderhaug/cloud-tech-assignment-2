package functions

import (
	"corona-information-service/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	return res
}

// UnmarshalResponse Method for unmarshalling GraphQL response into a struct */
func UnmarshalResponse(res *http.Response) model.Response {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response model.Response

	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal(err)
	}

	return response
}

// EncodeCaseInformation */
func EncodeCaseInformation(w http.ResponseWriter, r model.Case) {
	// Write content type header
	w.Header().Add("content-type", "application/json")

	// Instantiate encoder
	encoder := json.NewEncoder(w)

	//Encodes response
	err := encoder.Encode(r)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}
}
