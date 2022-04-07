package customhttp

import (
	"net/http"
	"time"
)

//HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client Declare the client
var (
	Client HTTPClient
)

//init the client
func init() {
	Client = &http.Client{
		Timeout: 10 * time.Second,
	}
}
