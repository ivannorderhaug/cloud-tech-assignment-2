package _case

import (
	"corona-information-service/internal/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var endpoint *httptest.Server

var client = http.Client{}

type Tests struct {
	name       string
	server     *httptest.Server
	response   *model.Case
	statusCode int
}

func TestMain(m *testing.M) {
	endpoint = httptest.NewServer(http.HandlerFunc(CaseHandler(&client)))
	m.Run()
}

//TestCaseHandler tests different request cases
func TestCaseHandler(t *testing.T) {

	tests := []Tests{
		{
			name: "basic-request-with-existing-country",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				//This is the nested structure that gets returned when doing a graphql request
				_, _ = w.Write([]byte("{\n  \"data\": {\n    \"country\": {\n      \"name\": \"Norway\",\n      \"mostRecent\": {\n        \"date\": \"2022-04-05\",\n        \"confirmed\": 1411550,\n        \"recovered\": 0,\n        \"deaths\": 2518,\n        \"growthRate\": 0.0010630821154695824\n      }\n    }\n  }\n}"))
			})),
			response: &model.Case{
				Country:        "Norway",
				Date:           "2022-04-05",
				ConfirmedCases: 1411550,
				Recovered:      0,
				Deaths:         2518,
				GrowthRate:     0.0010630821154695824,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "basic-request-with-non-existing-country ",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("{\n  \"errors\": [\n    {\n      \"message\": \"Couldn't find data from country test\",\n      \"locations\": [\n        {\n          \"line\": 2,\n          \"column\": 3\n        }\n      ],\n      \"path\": [\n        \"country\"\n      ],\n      \"extensions\": {\n        \"code\": \"INTERNAL_SERVER_ERROR\"\n      }\n    }\n  ],\n  \"data\": {\n    \"country\": null\n  }\n}"))
			})),
			response:   &model.Case{},
			statusCode: http.StatusNotFound,
		},
		{
			name:       "basic-request-with-no-parameters ",
			server:     nil,
			response:   nil,
			statusCode: http.StatusBadRequest,
		},
	}

	//Request with existing country name as parameter, expects populated case struct and status code 200
	t.Run(tests[0].name, func(t *testing.T) {
		model.CASES_URL = tests[0].server.URL

		request, _ := http.NewRequest(http.MethodGet, endpoint.URL+model.CASE_PATH+"Norway", nil)
		res, _ := client.Do(request)

		decoder := json.NewDecoder(res.Body)

		var c model.Case
		_ = decoder.Decode(&c)

		if !reflect.DeepEqual(tests[0].response, &c) {
			t.Errorf("FAILED: expected %v, got %v\n", tests[0].response, &c)
		}

		if !reflect.DeepEqual(tests[0].statusCode, res.StatusCode) {
			t.Errorf("FAILED: expected %v, got %v\n", tests[0].statusCode, res.StatusCode)
		}
	})

	//Request with non-existing country as parameter, expects empty case struct and status code 404
	t.Run(tests[1].name, func(t *testing.T) {
		model.CASES_URL = tests[1].server.URL

		request, _ := http.NewRequest(http.MethodGet, endpoint.URL+model.CASE_PATH+"Test", nil)
		res, _ := client.Do(request)

		decoder := json.NewDecoder(res.Body)

		var c model.Case
		_ = decoder.Decode(&c)

		if !reflect.DeepEqual(tests[1].response, &c) {
			t.Errorf("FAILED: expected %v, got %v\n", tests[1].response, &c)
		}

		if !reflect.DeepEqual(tests[1].statusCode, res.StatusCode) {
			t.Errorf("FAILED: expected %v, got %v\n", tests[1].statusCode, res.StatusCode)
		}
	})

	//Request with no parameters, expects status code 400
	t.Run(tests[2].name, func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, endpoint.URL+model.CASE_PATH, nil)
		res, _ := client.Do(request)

		if !reflect.DeepEqual(tests[2].statusCode, res.StatusCode) {
			t.Errorf("FAILED: expected %v, got %v\n", tests[2].statusCode, res.StatusCode)
		}
	})

}
