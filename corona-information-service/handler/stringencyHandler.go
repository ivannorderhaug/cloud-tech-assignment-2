package handler

import (
	"corona-information-service/functions"
	"corona-information-service/model"
	"encoding/json"
	"net/http"
	"strings"
)

// StringencyHandler */
func StringencyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}

	path, ok, msg := functions.PathSplitter(r.URL.Path, 1)
	if !ok {
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	//Country name or isocode.
	s := strings.ToUpper(path[0])
	if len(s) != 3 {
		http.Error(w, "Invalid alpha-3 country code. Please try again. ", http.StatusBadRequest)
		return
	}

}

// IssueRequest */
func IssueRequest(url string) (*http.Response, error) {
	// Create new request
	r, err := http.NewRequest(http.MethodGet, url, nil)
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

// DecodeResponse */
func DecodeResponse(res *http.Response) (model.CovidPolicyWrapper, error) {
	var w model.CovidPolicyWrapper

	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&w); err != nil {
		return model.CovidPolicyWrapper{}, err
	}

	return w, nil
}
