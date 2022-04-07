package customhttp

import (
	"bytes"
	"net/http"
)

// IssueRequest */
func IssueRequest(client HTTPClient, method string, url string, body []byte) (*http.Response, error) {
	// Create new request
	r, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return &http.Response{}, err
	}
	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		return &http.Response{}, err
	}

	return res, nil
}
