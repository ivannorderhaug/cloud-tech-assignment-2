package model

type Case struct {
	Country        string  `json:"country"`
	Date           string  `json:"date"`
	ConfirmedCases int     `json:"confirmed"`
	Recovered      int     `json:"recovered"`
	Deaths         int     `json:"deaths"`
	GrowthRate     float64 `json:"growth_rate"`
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
	Webhooks      int    `json:"webhooks"`
	Version       string `json:"version"`
	Uptime        int    `json:"uptime"`
}

type Webhook struct {
	Invoked     string `json:"invoked_at,omitempty" firestore:"-"`
	ID          string `json:"id" firestore:"id"`
	Url         string `json:"url" firestore:"url"`
	Country     string `json:"country" firestore:"country"`
	Calls       int    `json:"calls" firestore:"calls"`
	ActualCalls int    `json:"-" firestore:"actual_calls"`
}
