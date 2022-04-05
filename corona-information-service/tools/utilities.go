package tools

import (
	"regexp"
	"strings"
)

//PathSplitter
//Example of usage: /corona/v1/cases/norway has a length of 5, if it matches basePathLength(which is 4)+length param(which is 1 in this case).
//It'll return a slice containing only 1 element which is the search param (norway)*/
func PathSplitter(path string, length int) ([]string, bool) {
	//Trims away "/" at the end of path. Only if there is one there
	path = strings.TrimSuffix(path, "/")
	//Splits the path into a slice, separating each part by "/"
	parts := strings.Split(path, "/")
	//Gets the length of the basePath. Length will be 4.
	basePathLength := len(strings.Split("/corona/v1/", "/"))

	if len(parts) == basePathLength {
		//Returns empty slice with an error message as the path didn't match the required format
		return []string{}, false
	}

	//Compares length of parts slice with basePath length+length param
	if len(parts) != basePathLength+length {
		//Returns empty slice with an error message as the path didn't match the required format
		return []string{}, false
	}
	return parts[basePathLength : basePathLength+length], true
}

//IsValidDate Uses Regular Expressions to validate if string matches required format */
func IsValidDate(date string) bool {
	//YYYY-mm-dd
	pattern := regexp.MustCompile("([12]\\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01]))")
	return pattern.MatchString(date)
}
