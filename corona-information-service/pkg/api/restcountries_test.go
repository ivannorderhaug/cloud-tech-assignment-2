package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type Test struct {
	name     string
	server   *httptest.Server
	response interface{}
}

func TestGetCountryNameByAlphaCode(t *testing.T) {
	test := []Test{{
		name: "basic-request-with-real-alpha3",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{\"name\":{\"common\":\"Norway\",\"official\":\"Kingdom of Norway\",\"nativeName\":{\"nno\":{\"official\":\"Kongeriket Noreg\",\"common\":\"Noreg\"},\"nob\":{\"official\":\"Kongeriket Norge\",\"common\":\"Norge\"},\"smi\":{\"official\":\"Norgga gonagasriika\",\"common\":\"Norgga\"}}}}"))
		})),
		response: "Norway",
	},
		{
			name: "basic-request-with-fake-alpha3",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{\"status\":404,\"message\":\"Not Found\"}"))
			})),
			response: nil,
		},
	}

	//Request with existing alpha3 code
	t.Run(test[0].name, func(t *testing.T) {
		defer test[0].server.Close()
		RESTCOUNTRIES = test[0].server.URL + "/v3.1/alpha/%s?fields=name"

		fmt.Printf("MOCKED URL: %s\n", RESTCOUNTRIES)
		name, err := GetCountryNameByAlphaCode("Nor")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("GOT: " + fmt.Sprint(name))
		if !reflect.DeepEqual(test[0].response, name) {
			t.Errorf("FAILED: expected %v, got %v\n", test[0].response, name)
		}
	})

	//Request with non-existing alpha3 code
	t.Run(test[1].name, func(t *testing.T) {
		defer test[1].server.Close()
		RESTCOUNTRIES = test[1].server.URL + "/v3.1/alpha/%s?fields=name"

		fmt.Printf("MOCKED URL: %s\n", RESTCOUNTRIES)
		name, err := GetCountryNameByAlphaCode("test")
		if err != nil {
			fmt.Printf("GOT: %s \n", err)
		}
		if !reflect.DeepEqual(test[1].response, name) {
			t.Errorf("FAILED: expected %v, got %v\n", test[1].response, name)
		}
	})
}
