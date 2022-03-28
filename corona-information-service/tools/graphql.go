package tools

import (
	"encoding/json"
	"fmt"
)

func ConvertToGraphql(country string, requestBody string) ([]byte, error) {
	//Formats query
	query := fmt.Sprintf(requestBody, country)

	type GraphQLRequest struct {
		Query string `json:"query"`
	}

	jsonQuery, err := json.Marshal(GraphQLRequest{Query: query})
	if err != nil {
		return nil, err
	}
	return jsonQuery, nil
}
