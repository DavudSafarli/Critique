package testing_utils

import (
	"net/url"
)

// Get connection string for LOCAL TEST DATABASE
func GetTestDbConnStr() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("local-test-user", "local-test-pass"),
		Host:     "localhost:5433",
		Path:     "critique-local-test",
		RawQuery: "sslmode=disable",
	}
	connStr := u.String()
	return connStr
}
