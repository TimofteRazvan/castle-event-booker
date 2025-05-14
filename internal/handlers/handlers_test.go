package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/TimofteRazvan/castle-event-booker/internal/driver"
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

func TestNewRepo(t *testing.T) {
	fakeSQL := &sql.DB{}
	fakeDriver := &driver.DB{
		SQL: fakeSQL,
	}
	NewRepo(&app, fakeDriver)
}

func TestRepository_PostBooking(t *testing.T) {
	
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
		t.Errorf("MakeReservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
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
		t.Errorf("MakeReservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
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
		t.Errorf("MakeReservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostMakeReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "Knights' Hall",
		},
	}
	requestBody := "first_name=Mick"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last_name=Jagger")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=mick@gmail.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=123")
	request, err := http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context := getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()
	session.Put(context, "reservation", reservation)

	handler := http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	// test with reservation not in session
	request, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// test with missing form body
	request, err = http.NewRequest("POST", "/make-reservation", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	session.Put(context, "reservation", reservation)

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// test with invalid form data
	requestBody = "first_name=Mick"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last_name=Jagger")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=mick@com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=123")
	request, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	session.Put(context, "reservation", reservation)

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// test with invalid insert reservation
	requestBody = "first_name=Mick"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last_name=Jagger")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=mick@gmail.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=123")
	request, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reservation.RoomID = 99
	session.Put(context, "reservation", reservation)

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// test with invalid insert room restriction
	requestBody = "first_name=Mick"
	requestBody = fmt.Sprintf("%s&%s", requestBody, "last_name=Jagger")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "email=mick@gmail.com")
	requestBody = fmt.Sprintf("%s&%s", requestBody, "phone=123")
	request, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reservation.RoomID = 100
	session.Put(context, "reservation", reservation)

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func getContext(r *http.Request) context.Context {
	context, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return context
}
