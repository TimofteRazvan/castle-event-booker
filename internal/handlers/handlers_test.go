package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/TimofteRazvan/castle-event-booker/internal/models"
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
	// {"home", "/", "GET", []postData{}, http.StatusOK},
	// {"about", "/about", "GET", []postData{}, http.StatusOK},
	// {"knights", "/knights", "GET", []postData{}, http.StatusOK},
	// {"throne", "/throne", "GET", []postData{}, http.StatusOK},
	// {"banquet", "/banquet", "GET", []postData{}, http.StatusOK},
	// {"contact", "/contact", "GET", []postData{}, http.StatusOK},
	// {"booking", "/booking", "GET", []postData{}, http.StatusOK},
	// {"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// {"reservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	// {"post-booking", "/booking", "POST", []postData{
	// 	{key: "start", value: "2022-05-26"},
	// 	{key: "end", value: "2022-05-27"},
	// }, http.StatusOK},
	// {"post-booking-json", "/booking-json", "POST", []postData{
	// 	{key: "start", value: "2022-05-26"},
	// 	{key: "end", value: "2022-05-27"},
	// }, http.StatusOK},
	// {"post-make-reservation", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "Razvan"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "razvan@gmail.com"},
	// 	{key: "phone", value: "0777777777"},
	// }, http.StatusOK},
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
			values := url.Values{}
			for _, param := range test.params {
				values.Add(param.key, param.value)
			}

			response, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", test.name, test.expectedStatusCode, response.StatusCode)
			}
		}
	}
}

func TestRepository_Booking(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "Knights' Hall",
		},
	}

	request, err := http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}

	context := getContext(request)
	request = request.WithContext(context)

	responseRecorder := httptest.NewRecorder()
	session.Put(context, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Booking)
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}
}

func getContext(r *http.Request) context.Context {
	context, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return context
}
