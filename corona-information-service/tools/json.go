package tools

import (
	"encoding/json"
	"net/http"
)

// Encode */
func Encode(w http.ResponseWriter, data interface{}) {
	// Write content type header
	w.Header().Add("content-type", "application/json")

	// Instantiate encoder
	encoder := json.NewEncoder(w)

	//Encodes response
	err := encoder.Encode(&data)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}
}

// Decode */
func Decode(res *http.Response, data interface{}) error {
	dec := json.NewDecoder(res.Body)

	if err := dec.Decode(data); err != nil {
		return err
	}

	return nil
}
