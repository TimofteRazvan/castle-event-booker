package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
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
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"knights", "/knights", "GET", http.StatusOK},
	{"throne", "/throne", "GET", http.StatusOK},
	{"banquet", "/banquet", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"booking", "/booking", "GET", http.StatusOK},
	// {"reservation-summary", "/reservation-summary", "GET", http.StatusOK},
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
		response, err := testServer.Client().Get(testServer.URL + test.url)

		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if response.StatusCode != test.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", test.name, test.expectedStatusCode, response.StatusCode)
		}
	}
}

func TestRepository_MakeReservation(t *testing.T) {
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

	handler := http.HandlerFunc(Repo.MakeReservation)
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	// test with reservation not in session
	request, err = http.NewRequest("GET", "make-reservation", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// test with inexistent room
	request, err = http.NewRequest("GET", "make-reservation", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	responseRecorder = httptest.NewRecorder()
	reservation.RoomID = 99
	session.Put(context, "reservation", reservation)

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func getContext(r *http.Request) context.Context {
	context, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return context
}
