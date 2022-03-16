package model

type GraphQLRequest struct {
	Query string `json:"query"`
}

type Case struct {
	Name           string  `json:"country"`
	Date           string  `json:"date"`
	ConfirmedCases int     `json:"confirmed"`
	Recovered      int     `json:"recovered"`
	Deaths         int     `json:"deaths"`
	Growth         float64 `json:"growth_rate"`
}
