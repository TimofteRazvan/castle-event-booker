package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

var tests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"knights", "/knights", "GET", []postData{}, http.StatusOK},
	{"throne", "/throne", "GET", []postData{}, http.StatusOK},
	{"banquet", "/banquet", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"booking", "/booking", "GET", []postData{}, http.StatusOK},
	{"makeReservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"reservationSummary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range tests {
		if test.method == "GET" {
			response, err := testServer.Client().Get(testServer.URL + test.url)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", test.name, test.expectedStatusCode, response.StatusCode)
			}
		} else {

		}
	}
}
