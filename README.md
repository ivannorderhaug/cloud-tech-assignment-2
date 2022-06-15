# ASSIGNMENT 2 - CORONA INFORMATION SERVICE


### Improved API can be found here: https://github.com/ivannorderhaug/cloud-tech-assignment-3

## Assumptions

*  Active policies: in issue [#41](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2022/-/issues/41) it is stated that: 
`The response is the simple variant - the number of active policies returned, nothing else.` I interpret this as all the policies that gets returned, count as active. Even though some of them have "no measures" as value in the `policy_value_display_field`

## Overview

### Project structure
```
.
└── corona-information-service/
    ├── cmd/
    │   └── server.go
    └── internal/
        ├── handler/
        │   ├── cases/
        │   │   ├── caseHandler.go
        │   │   └── caseHandler_test.go
        │   ├── policy/
        │   │   ├── policyHandler.go
        │   │   └── policyHandler_test.go
        │   ├── deafult.go
        │   ├── notificationHandler.go
        │   └── status.go
        ├── model/
        │   ├── constants.go
        │   └── structs.go
        ├── pkg/
        │   ├── api/
        │   │   ├── restcountries.go
        │   │   └── restcountries_test.go
        │   ├── cache/
        │   │   └── cache.go
        │   ├── customhttp/
        │   │   ├── HTTPClient.go
        │   │   └── request.go
        │   ├── customjson/
        │   │   └── json.go
        │   ├── db/
        │   │   └── firestore.go
        │   └── webhook/
        │       └── webhooks.go
        ├── tools/
        │   ├── graphql/
        │   │   └── graphql.go
        │   ├── hash/
        │   │   └── hash.go
        │   └── utilities
        └── go.mod/
            └── go.sum
```

## Deployment

### Local service
There a two ways to deploy this service locally. Either by cloning the repository  or downloading the .zip file.  

```bash
git clone https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2022-workspace/ivann/assignment-2-server.git
```

[DOWNLOAD](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2022-workspace/ivann/assignment-2-server/-/archive/main/assignment-2-server-main.zip)

#### After installation/cloning: ```simply run ../../cmd/server.go ```

#### NOTE: By deploying the service locally, you lose access to the notification endpoint and the webhooks. To get access to them, please register [here](https://firebase.google.com/). Create a service account and put the .json file in the project folder --> (./corona-information-service) 

### Hosted service
There is a currently a running docker container for the service available at [10.212.136.78](http://10.212.136.78/). To access it, you must first connect to the NTNU VPN. Thereafter, you are free to use the service as you please.

## Endpoints
There are currently four endpoints available with the following resource paths:
```
/corona/v1/cases/
/corona/v1/policy/
/corona/v1/status/
/corona/v1/notifications/
```

### Covid-19 Cases per Country
The initial endpoint focuses on return the latest number of confirmed cases and deaths for a given country, alongside growth rate of cases.

#### - Request
```
Method: GET
Path: /corona/v1/cases/{:country_name}
```
```{:country_name}``` refers to the name for the country as supported by the Covid 19 cases API or the ISO 3166-1 alpha-3 country code.

Example requests:  
```/corona/v1/cases/Norway```   
```/corona/v1/cases/NOR```

#### - Response

- Content type: ```application/json```

Body (Example):
```json
{
    "country": "Norway",
    "date": "2022-03-05",
    "confirmed": 1305006,
    "recovered": 0,
    "deaths": 1664,
    "growth_rate": 0.004199149089414866
}
```
#### Countries covered by the case endpoint. 
The external API is extremely sensitive so make sure your search is correctly formatted.  
CTRL+F to see if your desired country is covered.
```json
    "Bangladesh",
    "Dominican Republic",
    "Gambia",
    "Guyana",
    "Mauritius",
    "Monaco",
    "Algeria",
    "Austria",
    "Ukraine",
    "Comoros",
    "Congo (Kinshasa)",
    "Holy See",
    "Iran",
    "Saudi Arabia",
    "Winter Olympics 2022",
    "Antigua and Barbuda",
    "Belgium",
    "Tanzania",
    "Turkey",
    "Vietnam",
    "Armenia",
    "Kiribati",
    "Maldives",
    "Nepal",
    "Egypt",
    "Jamaica",
    "France",
    "Korea, South",
    "Latvia",
    "Papua New Guinea",
    "US",
    "Barbados",
    "Ethiopia",
    "Central African Republic",
    "Cyprus",
    "Greece",
    "Guatemala",
    "Honduras",
    "Japan",
    "Albania",
    "Brazil",
    "Togo",
    "Namibia",
    "Pakistan",
    "Serbia",
    "Singapore",
    "Lebanon",
    "Liberia",
    "Zimbabwe",
    "Colombia",
    "Trinidad and Tobago",
    "Kazakhstan",
    "Kuwait",
    "Luxembourg",
    "Netherlands",
    "Qatar",
    "Bulgaria",
    "Burma",
    "MS Zaandam",
    "Peru",
    "Saint Vincent and the Grenadines",
    "Afghanistan",
    "India",
    "Chad",
    "Diamond Princess",
    "El Salvador",
    "Belize",
    "Cambodia",
    "Laos",
    "Malaysia",
    "Syria",
    "Thailand",
    "Cabo Verde",
    "Eritrea",
    "Indonesia",
    "Ireland",
    "Jordan",
    "Venezuela",
    "Bhutan",
    "Burkina Faso",
    "China",
    "Congo (Brazzaville)",
    "Ecuador",
    "Guinea-Bissau",
    "Malta",
    "Panama",
    "Australia",
    "Benin",
    "Saint Lucia",
    "Cuba",
    "Sierra Leone",
    "Slovakia",
    "Gabon",
    "Sao Tome and Principe",
    "Finland",
    "Hungary",
    "Burundi",
    "Chile",
    "Montenegro",
    "Spain",
    "Costa Rica",
    "Estonia",
    "North Macedonia",
    "Poland",
    "Grenada",
    "Italy",
    "Mongolia",
    "Russia",
    "Senegal",
    "Suriname",
    "Switzerland",
    "Yemen",
    "Angola",
    "Liechtenstein",
    "Cote d'Ivoire",
    "Germany",
    "Madagascar",
    "Mexico",
    "Rwanda",
    "South Sudan",
    "Antarctica",
    "Bosnia and Herzegovina",
    "Uganda",
    "Vanuatu",
    "Georgia",
    "Kenya",
    "Libya",
    "Lithuania",
    "Malawi",
    "Saint Kitts and Nevis",
    "Botswana",
    "Fiji",
    "Solomon Islands",
    "Cameroon",
    "Canada",
    "New Zealand",
    "Portugal",
    "Samoa",
    "Sweden",
    "Andorra",
    "Argentina",
    "Timor-Leste",
    "Zambia",
    "Dominica",
    "Micronesia",
    "Nicaragua",
    "Nigeria",
    "Sri Lanka",
    "Bahrain",
    "Belarus",
    "Oman",
    "Taiwan*",
    "Uzbekistan",
    "Brunei",
    "Mozambique",
    "Palau",
    "Paraguay",
    "Slovenia",
    "Tajikistan",
    "Iraq",
    "Mauritania",
    "San Marino",
    "Seychelles",
    "Tunisia",
    "Uruguay",
    "Bolivia",
    "Iceland",
    "Romania",
    "South Africa",
    "West Bank and Gaza",
    "Denmark",
    "Eswatini",
    "Lesotho",
    "Guinea",
    "Mali",
    "Niger",
    "Philippines",
    "Somalia",
    "Summer Olympics 2020",
    "Bahamas",
    "Djibouti",
    "Tonga",
    "Ghana",
    "Haiti",
    "Moldova",
    "Morocco",
    "Croatia",
    "Czechia",
    "Israel",
    "Kosovo",
    "Kyrgyzstan",
    "Marshall Islands",
    "Norway",
    "Sudan",
    "Azerbaijan",
    "Equatorial Guinea",
    "United Arab Emirates",
    "United Kingdom"
```

### Covid Policy Stringency per Country
The second endpoint provides an overview of the current stringency level of policies regarding Covid-19 for a given country, in addition to the number of currently active policies.

#### - Request
```
Method: GET
Path: /corona/v1/policy/{:country_code}{?scope=YYYY-MM-DD}
```
```{:country_code}``` refers to the ISO 3166-1 alpha-3 country code.  
```{?scope=YYYY-MM-DD}```optional: indicates the date for which the policy stringency information should be returned.

#### - Response

- Content type: ```application/json```
Body (Example):
```json
{
    "country_code": "NOR",
    "scope": "2021-12-12",
    "stringency": 40.74,
    "policies": 20
}

```
### Status Interface
The status interface indicates the availability of all individual services this service depends on. The status interface further provides information about the number of registered webhooks (more details is provided in the next section), and the uptime of the service.

#### - Request
```
Method: GET
Path: /corona/v1/status
```

- Content type: ```application/json```
Body (Example):
```json
{
    "cases_api": "200 OK",
    "policy_api": "200 OK",
    "restcountries_api": "200 OK",
    "webhooks": 1,
    "version": "v1",
    "uptime": "6 s"
}
```
### Notification
As an additional feature, users can register webhooks that are triggered by the service based on specified events, specifically if information about given countries is invoked, where the minimum frequency can be specified. Users can register multiple webhooks. 
### Registration of Webhook
#### - Request
```
Method: POST
Path: /corona/v1/notifications/
```

- Content type: ```application/json```

The body contains  
- The URL to be triggered upon event (the service that should be invoked)
- The country for which the trigger applies. Supports usage of ISO 3166-1 alpha-3 country code. Country code will be converted to the appropriate country name.
- The minimum number of repeated invocations before notification should occur


Body (Example):
```json
{
   "url": "http://10.212.136.78:8081/client",
   "country": "Norway",
   "calls": 2
}
```
#### - Response

- Content type: ```application/json```   
 
Body (Example):
```json
{
    "webhook_id": "yzdPGfAORdmgKDv"
}
```
### Deletion of Webhook

#### - Request

```
Method: DELETE
Path: /corona/v1/notifications/{id}
```

```{id}``` is the ID returned during the webhook registration

#### - Response
- Content type: ```application/json```  
```json
{
    "result": "The webhook has been successfully removed from the database!"
}
```


### View registered webhook

#### - Request

```
Method: GET
Path: /corona/v1/notifications/{id}
```

```{id}``` is the ID returned during the webhook registration

#### - Response
* Content type: `application/json`

Body (Example):
```json
{
    "id": "yzdPGfAORdmgKDv",
    "url": "http://10.212.136.78:8081/client",
    "country": "Norway",
    "calls": 2
}

```

### View all registered webhooks

#### - Request

```
Method: GET
Path: /corona/v1/notifications/
```

#### - Response

The response is a collection of all registered webhooks.

* Content type: `application/json`

Body (Example):
```json
[{
    "id": "yzdPGfAORdmgKDv",
    "url": "http://10.212.136.78:8081/client",
    "country": "Norway",
    "calls": 2
 },
 {
    "webhook_id": "DiSoisivucios",
    "url": "http://10.212.136.78:8081/client",
    "country": "Sweden",
    "calls": 5
 },
...
]
```
### Webhook Invocation (upon trigger)

When a webhook is triggered, it sends information as follows.

```
Method: POST
Path: <url specified in the corresponding webhook registration>
```

* Content type: `application/json`

Body (Example):
```json
{
    "id": "yzdPGfAORdmgKDv",
    "url": "http://10.212.136.78:8081/client",
    "country": "Norway",
    "calls": 2
}
```
