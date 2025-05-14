package dbrepo

import (
	"time"

	"github.com/TimofteRazvan/castle-event-booker/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction inserts a room_restriction into the database
func (m *testDBRepo) InsertRoomRestriction(rr models.RoomRestriction) error {
	return nil
}

// SearchAvailabilityByDateByRoomID returns true if the room in the date interval is available for booking, false otherwise
func (m *testDBRepo) SearchAvailabilityByDateByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityByDateAllRooms returns a slice of available rooms for given date interval
func (m *testDBRepo) SearchAvailabilityByDateAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID returns the room with the given ID
func (m *testDBRepo) GetRoomByID(roomID int) (models.Room, error) {
	var room models.Room
	return room, nil
}
