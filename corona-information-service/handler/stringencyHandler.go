package handler

import (
	"corona-information-service/functions"
	"net/http"
	"strings"
)

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
