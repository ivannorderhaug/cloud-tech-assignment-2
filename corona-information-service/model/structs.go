package model

type GraphQLRequest struct {
	Query string `json:"query"`
}
type Info struct {
	Date       string  `json:"date"`
	Confirmed  int     `json:"confirmed"`
	Recovered  int     `json:"recovered"`
	Deaths     int     `json:"deaths"`
	GrowthRate float64 `json:"growthRate"`
}

type CountryInfo struct {
	Name string `json:"name"`
	Info Info   `json:"mostRecent"`
}

type Country struct {
	Country CountryInfo `json:"country"`
}

type Response struct {
	Data Country `json:"data"`
}

type Case struct {
	Country        string  `json:"country"`
	Date           string  `json:"date"`
	ConfirmedCases int     `json:"confirmed"`
	Recovered      int     `json:"recovered"`
	Deaths         int     `json:"deaths"`
	GrowthRate     float64 `json:"growth_rate"`
}
