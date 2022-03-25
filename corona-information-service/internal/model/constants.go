package model

//PATH
const DEFAULT_PATH = "/"
const CASE_PATH = "/corona/v1/cases/"
const POLICY_PATH = "/corona/v1/policy/"
const STATUS_PATH = "/corona/v1/status/"
const NOTIFICATION_PATH = "/corona/v1/notifications/"
const VERSION = "v1"

//URL
const CASES_URL = "https://covid19-graphql.vercel.app/"
const STRINGENCY_URL = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/actions/"

const QUERY = "query {\n  country(name: \"%s\") {\n    name\n    mostRecent {\n      date(format: \"yyyy-MM-dd\")\n      confirmed\n      recovered\n      deaths\n      growthRate\n    }\n  }\n}"
