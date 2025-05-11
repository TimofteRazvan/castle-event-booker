package repository

import "github.com/TimofteRazvan/castle-event-booker/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(rr models.RoomRestriction) error
}
