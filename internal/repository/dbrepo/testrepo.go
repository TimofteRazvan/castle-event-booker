package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/TimofteRazvan/castle-event-booker/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 99 {
		return 0, errors.New("reservation: inexistent room ID")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room_restriction into the database
func (m *testDBRepo) InsertRoomRestriction(rr models.RoomRestriction) error {
	if rr.RoomID == 100 {
		return errors.New("room restriction: inexistent room ID")
	}
	return nil
}

// SearchAvailabilityByDateByRoomID returns true if the room in the date interval is available for booking, false otherwise
func (m *testDBRepo) SearchAvailabilityByDateByRoomID(start, end time.Time, roomID int) (bool, error) {
	layout := "2006-01-02"
	// if the start date is after 2060-12-31, then return empty slice,
	// indicating no rooms are available;
	str := "2060-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// test to fail the query -- 2060-01-01 as start
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return false, errors.New("some error")
	}

	if start.After(t) {
		return false, nil
	}

	return true, nil
}

// SearchAvailabilityByDateAllRooms returns a slice of available rooms for given date interval
func (m *testDBRepo) SearchAvailabilityByDateAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room

	// if the start date is after 2060-12-31, then return empty slice,
	// indicating no rooms are available;
	layout := "2006-01-02"
	str := "2060-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return rooms, errors.New("some error")
	}

	if start.After(t) {
		return rooms, nil
	}

	// some room is available for date interval
	room := models.Room{
		ID: 1,
	}
	rooms = append(rooms, room)

	return rooms, nil
}

// GetRoomByID returns the room with the given ID
func (m *testDBRepo) GetRoomByID(roomID int) (models.Room, error) {
	var room models.Room
	if roomID < 2 || roomID > 4 {
		return room, errors.New("inexistent room ID")
	}

	return room, nil
}

func (m *testDBRepo) GetUserByID(userID int) (models.User, error) {
	return models.User{}, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 0, "", nil
}

// AllReservations returns a slice of all reservations
func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}
