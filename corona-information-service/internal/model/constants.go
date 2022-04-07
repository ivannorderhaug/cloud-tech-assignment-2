package model

// PATHS

const DEFAULT_PATH = "/"
const CASE_PATH = "/corona/v1/cases/"
const POLICY_PATH = "/corona/v1/policy/"
const STATUS_PATH = "/corona/v1/status/"
const NOTIFICATION_PATH = "/corona/v1/notifications/"
const VERSION = "v1"

// URLS TO INVOKE

var CASES_URL = "https://covid19-graphql.vercel.app/"
var STRINGENCY_URL = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/"

// GRAPHQL

const QUERY = "query {\n  country(name: \"%s\") {\n    name\n    mostRecent {\n      date(format: \"yyyy-MM-dd\")\n      confirmed\n      recovered\n      deaths\n      growthRate\n    }\n  }\n}"

// USED TO CHECK STATUS OF API

const RESTCOUNTRIES_API = "https://restcountries.com/v3.1/all"
const STRINGENCY_API = "https://covidtrackerapi.bsg.ox.ac.uk/api"
const CASES_API = "https://covid19-graphql.vercel.app/?query=%7B__typename%7D"
