package _case

import (
	"corona-information-service/internal/model"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type Tests struct {
	name          string
	server        *httptest.Server
	response      *model.Case
	expectedError error
}

func TestGetCase(t *testing.T) {
	tests := []Tests{
		{
			name: "existing-country-name",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				//This is the nested structure that gets returned when doing a graphql request
				w.Write([]byte("{\n  \"data\": {\n    \"country\": {\n      \"name\": \"Norway\",\n      \"mostRecent\": {\n        \"date\": \"2022-04-05\",\n        \"confirmed\": 1411550,\n        \"recovered\": 0,\n        \"deaths\": 2518,\n        \"growthRate\": 0.0010630821154695824\n      }\n    }\n  }\n}"))
			})),
			response: &model.Case{
				Country:        "Norway",
				Date:           "2022-04-05",
				ConfirmedCases: 1411550,
				Recovered:      0,
				Deaths:         2518,
				GrowthRate:     0.0010630821154695824,
			},
			expectedError: nil,
		},
		{
			name: "non-existing-country-name",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				//This is the nested structure that gets returned when doing a graphql request
				w.Write([]byte("{\n  \"errors\": [\n    {\n      \"message\": \"Couldn't find data from country test\",\n      \"locations\": [\n        {\n          \"line\": 2,\n          \"column\": 3\n        }\n      ],\n      \"path\": [\n        \"country\"\n      ],\n      \"extensions\": {\n        \"code\": \"INTERNAL_SERVER_ERROR\"\n      }\n    }\n  ],\n  \"data\": {\n    \"country\": null\n  }\n}"))
			})),
			response:      nil,
			expectedError: nil,
		},
	}

	t.Run(tests[0].name, func(t *testing.T) {
		defer tests[0].server.Client()

		res, err := getCase(tests[0].server.URL, "Norway")

		if !reflect.DeepEqual(res, tests[0].response) {
			t.Errorf("FAILED: expected %v, got %v\n", tests[0].response, res)
		}
		if !errors.Is(err, tests[0].expectedError) {
			t.Errorf("FAILED: expected %v, got %v\n", tests[0].expectedError, err)
		}
	})

	t.Run(tests[1].name, func(t *testing.T) {
		defer tests[1].server.Client()

		res, err := getCase(tests[1].server.URL, "Norway")

		if !reflect.DeepEqual(res, tests[1].response) {
			t.Errorf("FAILED: expected %v, got %v\n", tests[1].response, res)
		}
		if !errors.Is(err, tests[1].expectedError) {
			t.Errorf("FAILED: expected %v, got %v\n", tests[1].expectedError.Error(), err.Error())
		}
	})

}
