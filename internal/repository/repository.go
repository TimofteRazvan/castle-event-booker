package repository

import (
	"time"

	"github.com/TimofteRazvan/castle-event-booker/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(rr models.RoomRestriction) error

	SearchAvailabilityByDate(start, end time.Time, roomID int) (bool, error)
}
