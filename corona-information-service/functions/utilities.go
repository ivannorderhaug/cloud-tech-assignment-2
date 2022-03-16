package functions

import (
	"strings"
)

//PathSplitter
func PathSplitter(path string, length int) ([]string, bool, string) {
	//Removes last "/" of the url
	path = strings.TrimSuffix(path, "/")

	//Splits the url into an slice using "/" as a separator
	parts := strings.Split(path, "/")

	//Determines the length of the base path of the API.
	basePathLength := len(strings.Split("/corona/v1/", "/"))

	//Compares length of full path slice and length of base path slice + length param.
	if len(parts) != basePathLength+length {
		return []string{}, false, "Does not match expected path format "
	}

	//Returns a part of the parts slice without the base path parts
	return parts[basePathLength : basePathLength+length], true, ""
}
