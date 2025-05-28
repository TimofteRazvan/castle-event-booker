package repository

import (
	"time"

	"github.com/TimofteRazvan/castle-event-booker/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(rr models.RoomRestriction) error

	SearchAvailabilityByDateByRoomID(start, end time.Time, roomID int) (bool, error)

	SearchAvailabilityByDateAllRooms(start, end time.Time) ([]models.Room, error)

	GetRoomByID(roomID int) (models.Room, error)

	GetUserByID(userID int) (models.User, error)

	UpdateUser(u models.User) error

	Authenticate(email, testPassword string) (int, string, error)

	AllReservations() ([]models.Reservation, error)
}
