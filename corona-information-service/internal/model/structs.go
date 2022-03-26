package model

type GraphQLRequest struct {
	Query string `json:"query"`
}

// TmpCase Used to unwrap nested structure
type TmpCase struct {
	Data struct {
		Country struct {
			Name       string `json:"name"`
			MostRecent struct {
				Date       string  `json:"date"`
				Confirmed  int     `json:"confirmed"`
				Recovered  int     `json:"recovered"`
				Deaths     int     `json:"deaths"`
				GrowthRate float64 `json:"growthRate"`
			} `json:"mostRecent"`
		} `json:"country"`
	} `json:"data"`
}

type Case struct {
	Country        string  `json:"country"`
	Date           string  `json:"date"`
	ConfirmedCases int     `json:"confirmed"`
	Recovered      int     `json:"recovered"`
	Deaths         int     `json:"deaths"`
	GrowthRate     float64 `json:"growth_rate"`
}

// TmpPolicy Used to unwrap nested structure
type TmpPolicy struct {
	StringencyData struct {
		Stringency       float64 `json:"stringency"`
		StringencyActual float64 `json:"stringency_actual,omitempty"`
	} `json:"stringencyData"`
	PolicyActions []interface{} `json:"policyActions"`
}

type Policy struct {
	CountryCode string      `json:"country_code"`
	Scope       string      `json:"scope"`
	Stringency  float64     `json:"stringency,omitempty"`
	Policies    interface{} `json:"policies,omitempty"`
}

type Status struct {
	CasesApi      string `json:"cases_api"`
	PolicyApi     string `json:"policy_api"`
	RestCountries string `json:"restcountries_api"`
	//TODO Add webhooks
	Version string `json:"version"`
	Uptime  int    `json:"uptime"`
}

type Webhook struct {
	ID          string `json:"id,omitempty" firestore:"-"`
	Url         string `json:"url" firestore:"url"`
	Country     string `json:"country" firestore:"country"`
	Calls       int    `json:"calls" firestore:"calls"`
	ActualCalls int    `json:"-" firestore:"actual_calls"`
}
