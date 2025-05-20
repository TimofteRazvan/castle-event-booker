package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/TimofteRazvan/castle-event-booker/internal/driver"
	"github.com/TimofteRazvan/castle-event-booker/internal/models"
)

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
	testRepo := NewRepo(&app, fakeDriver)
	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

func TestRepository_PostBooking(t *testing.T) {
	// Case 1: all is correct + available rooms
	postData := url.Values{}
	postData.Add("start", "2025-07-01")
	postData.Add("end", "2025-07-02")
	request, err := http.NewRequest("POST", "/booking", strings.NewReader(postData.Encode()))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context := getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostBooking)
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("PostBooking handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	// Case 2: all is correct + no available rooms
	postData = url.Values{}
	postData.Add("start", "2070-07-01")
	postData.Add("end", "2070-07-02")
	request, err = http.NewRequest("POST", "/booking", strings.NewReader(postData.Encode()))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("PostBooking handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	// Case 3: missing form body
	request, err = http.NewRequest("POST", "/booking", nil)
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
		t.Errorf("PostBooking handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Case 4: invalid start date format
	postData = url.Values{}
	postData.Add("start", "20250701")
	postData.Add("end", "2025-07-02")
	request, err = http.NewRequest("POST", "/booking", strings.NewReader(postData.Encode()))
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
		t.Errorf("PostBooking handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Case 5: invalid end date format
	postData = url.Values{}
	postData.Add("start", "2025-07-01")
	postData.Add("end", "20250702")
	request, err = http.NewRequest("POST", "/booking", strings.NewReader(postData.Encode()))
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
		t.Errorf("PostBooking handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Case 6: failed searching due to database query restrictions
	postData = url.Values{}
	postData.Add("start", "2060-01-01")
	postData.Add("end", "2025-07-02")
	request, err = http.NewRequest("POST", "/booking", strings.NewReader(postData.Encode()))
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
		t.Errorf("PostBooking handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	// Case 1: all is correct
	reservation := models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "Knights' Hall",
		},
	}

	request, err := http.NewRequest("GET", "/choose-room", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context := getContext(request)
	request = request.WithContext(context)
	request.RequestURI = "/choose-room/2"

	responseRecorder := httptest.NewRecorder()
	session.Put(context, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	// Case 2: missing url parameter
	request, err = http.NewRequest("GET", "/choose-room", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.RequestURI = "/choose-room/badID"

	responseRecorder = httptest.NewRecorder()
	session.Put(context, "reservation", reservation)

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Case 3: reservation not in session
	request, err = http.NewRequest("GET", "/choose-room", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.RequestURI = "/choose-room/2"

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_BookingJSON(t *testing.T) {
	// Case 1: all is correct
	postData := url.Values{}
	postData.Add("start", "2025-07-01")
	postData.Add("end", "2025-07-02")
	postData.Add("room_id", "2")

	request, err := http.NewRequest("POST", "/booking-json", strings.NewReader(postData.Encode()))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context := getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.BookingJSON)
	handler.ServeHTTP(responseRecorder, request)

	var jsonResp jsonResponse
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Error("failed to parse json")
	}

	if !jsonResp.OK {
		t.Error("Room unavailable in BookingJSON, should be available")
	}

	// Case 2: room unavailable
	postData = url.Values{}
	postData.Add("start", "2070-01-01")
	postData.Add("end", "2070-01-02")
	postData.Add("room_id", "2")

	request, err = http.NewRequest("POST", "/booking-json", strings.NewReader(postData.Encode()))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Error("failed to parse json")
	}

	if jsonResp.OK {
		t.Error("Room available in BookingJSON, should be unavailable")
	}

	// Case 3: no request body
	request, err = http.NewRequest("POST", "/booking-json", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Error("failed to parse json")
	}

	if jsonResp.OK || jsonResp.Message != "Internal server error" {
		t.Error("Room is available but request is body empty")
	}

	// Case 4: fail parsing start date
	postData = url.Values{}
	postData.Add("start", "wrong")
	postData.Add("end", "2025-07-02")
	postData.Add("room_id", "2")

	request, err = http.NewRequest("POST", "/booking-json", strings.NewReader(postData.Encode()))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Error("failed to parse json")
	}

	if jsonResp.OK || jsonResp.Message != "Error parsing start date" {
		t.Error("Failed to parse start date in BookingJSON")
	}

	// Case 5: fail parsing end date
	postData = url.Values{}
	postData.Add("start", "2025-09-09")
	postData.Add("end", "wrong")
	postData.Add("room_id", "2")

	request, err = http.NewRequest("POST", "/booking-json", strings.NewReader(postData.Encode()))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Error("failed to parse json")
	}

	if jsonResp.OK || jsonResp.Message != "Error parsing end date" {
		t.Error("Failed to parse end date in BookingJSON")
	}

	// Case 6: fail parsing room id
	postData = url.Values{}
	postData.Add("start", "2025-09-09")
	postData.Add("end", "2025-09-10")
	postData.Add("room_id", "wrong")

	request, err = http.NewRequest("POST", "/booking-json", strings.NewReader(postData.Encode()))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Error("failed to parse json")
	}

	if jsonResp.OK || jsonResp.Message != "Error converting string to int" {
		t.Error("Failed converting room id in BookingJSON")
	}

	// Case 7: database query error
	postData = url.Values{}
	postData.Add("start", "2060-01-01")
	postData.Add("end", "2060-01-02")
	postData.Add("room_id", "2")

	request, err = http.NewRequest("POST", "/booking-json", strings.NewReader(postData.Encode()))
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	err = json.Unmarshal(responseRecorder.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Error("failed to parse json")
	}

	if jsonResp.OK || jsonResp.Message != "Error querying database" {
		t.Error("Room is available but request is body empty")
	}
}

func TestRepository_BookRoom(t *testing.T) {
	// Case 1: all is correct
	reservation := models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "Knights' Hall",
		},
	}

	request, err := http.NewRequest("GET", "/book-room?s=2025-09-09&e=2025-09-02&id=2", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context := getContext(request)
	request = request.WithContext(context)
	responseRecorder := httptest.NewRecorder()
	session.Put(context, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusSeeOther)
	}

	// Case 2: start date parse failed
	request, err = http.NewRequest("GET", "/book-room?s=wrong&e=2025-09-02&id=2", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	responseRecorder = httptest.NewRecorder()
	session.Put(context, "reservation", reservation)

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Case 3: end date parse failed
	request, err = http.NewRequest("GET", "/book-room?s=2025-09-01&e=wrong&id=2", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	responseRecorder = httptest.NewRecorder()
	session.Put(context, "reservation", reservation)

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Case 4: id parse failed
	request, err = http.NewRequest("GET", "/book-room?s=2025-09-01&e=2025-09-02&id=wrong", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	responseRecorder = httptest.NewRecorder()
	session.Put(context, "reservation", reservation)

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}

	// Case 5: database id query failed
	request, err = http.NewRequest("GET", "/book-room?s=2025-09-01&e=2025-09-02&id=99", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	responseRecorder = httptest.NewRecorder()
	session.Put(context, "reservation", reservation)

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_MakeReservation(t *testing.T) {
	// Case 1: all is correct
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

	// Case 2: reservation not in session
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

	// Case 3: inexistent room
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
	// Case 1: all is correct
	reservation := models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "Knights' Hall",
		},
	}
	postData := url.Values{}
	postData.Add("first_name", "Mick")
	postData.Add("last_name", "Jagger")
	postData.Add("email", "mick@gmail.com")
	postData.Add("phone", "123")

	request, err := http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
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

	// Case 2: reservation not in session
	request, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
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

	// Case 3: missing form body
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

	// Case 4: invalid form data
	postData = url.Values{}
	postData.Add("first_name", "Mick")
	postData.Add("last_name", "Jagger")
	postData.Add("email", "mick@com")
	postData.Add("phone", "123")
	request, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
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

	// Case 5: invalid insert reservation room id
	postData = url.Values{}
	postData.Add("first_name", "Mick")
	postData.Add("last_name", "Jagger")
	postData.Add("email", "mick@gmail.com")
	postData.Add("phone", "123")
	request, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
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

	// Case 6: invalid insert room restriction
	postData = url.Values{}
	postData.Add("first_name", "Mick")
	postData.Add("last_name", "Jagger")
	postData.Add("email", "mick@gmail.com")
	postData.Add("phone", "123")
	request, err = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
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

func TestRepository_ReservationSummary(t *testing.T) {
	// Case 1: everything is correct
	reservation := models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "Knights' Hall",
		},
	}

	request, err := http.NewRequest("GET", "/reservation-summary", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context := getContext(request)
	request = request.WithContext(context)
	request.RequestURI = "/reservation-summary"

	responseRecorder := httptest.NewRecorder()
	session.Put(context, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusOK)
	}

	// Case 2: reservation not in session
	request, err = http.NewRequest("GET", "/reservation-summary", nil)
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}
	context = getContext(request)
	request = request.WithContext(context)
	request.RequestURI = "/reservation-summary"

	responseRecorder = httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationSummary handler returned wrong response code: got %d, wanted %d", responseRecorder.Code, http.StatusTemporaryRedirect)
	}
}

func getContext(r *http.Request) context.Context {
	context, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return context
}
