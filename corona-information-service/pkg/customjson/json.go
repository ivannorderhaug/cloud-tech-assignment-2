package customjson

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
func Decode(i interface{}, data interface{}) error {
	//declare the decodes
	var dec *json.Decoder

	//create decoder based on interface type
	switch i.(type) {
	case *http.Request:
		dec = json.NewDecoder(i.(*http.Request).Body)
	case *http.Response:
		dec = json.NewDecoder(i.(*http.Response).Body)
	}

	//decodes interface
	if err := dec.Decode(data); err != nil {
		return err
	}

	return nil
}
