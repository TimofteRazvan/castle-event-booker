package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/TimofteRazvan/castle-event-booker/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `INSERT INTO reservations(first_name, last_name, email, phone, 
			start_date, end_date, room_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room_restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(rr models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions(start_date, end_date, 
	room_id, reservation_id, restriction_id, created_at, updated_at)
	values($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		rr.StartDate,
		rr.EndDate,
		rr.RoomID,
		rr.ReservationID,
		rr.RestrictionID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDateByRoomID returns true if the room in the date interval is available for booking, false otherwise
func (m *postgresDBRepo) SearchAvailabilityByDateByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id)
			from room_restrictions rr 
			where room_id = $1 
			and $2 < end_date 
			and $3 > start_date`

	var numRows int
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityByDateAllRooms returns a slice of available rooms for given date interval
func (m *postgresDBRepo) SearchAvailabilityByDateAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select r.id, r.room_name
			from rooms r
			where r.id not in (
				select room_id
				from room_restrictions rr
				where $1 < rr.end_date and $2 > rr.start_date)`

	var rooms []models.Room
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRoomByID returns the room with the given ID
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, room_name, created_at, updated_at
			from rooms
			where id = $1`

	var room models.Room
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}

// GetUserByID returns the user with the given ID
func (m *postgresDBRepo) GetUserByID(userID int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at
			from users
			where id = $1`

	row := m.DB.QueryRowContext(ctx, query, userID)
	var u models.User
	err := row.Scan(&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt)
	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users
			set first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5
			where id=$6`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		u.UpdatedAt,
		u.ID)

	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user by checking the testPassword's hash against the stored hash
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	query := `select id, password
			from users
			where email=$1`

	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// AllReservations returns a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation
	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
	r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
			from reservations r
			left join rooms rm on r.room_id = rm.id
			order by r.start_date asc
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomID,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Processed,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// NewReservations returns a slice of new unprocessed reservations
func (m *postgresDBRepo) NewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
	r.end_date, r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name
			from reservations r
			left join rooms rm on r.room_id = rm.id
			where processed = 0
			order by r.start_date asc
			`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Email,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomID,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, reservation)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// GetReservationById returns a reservation via id
func (m *postgresDBRepo) GetReservationById(reservationID int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation models.Reservation

	query := `select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, 
	r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
			from reservations r
			left join rooms rm on r.room_id = rm.id
			where r.id = $1
			`

	row := m.DB.QueryRowContext(ctx, query, reservationID)
	err := row.Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomID,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Processed,
		&reservation.Room.ID,
		&reservation.Room.RoomName,
	)
	if err != nil {
		return reservation, err
	}

	return reservation, nil
}

// UpdateReservation updates a reservation in the database
func (m *postgresDBRepo) UpdateReservation(reservation models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations
			set first_name=$1, last_name=$2, email=$3, phone=$4, updated_at=$5
			where id=$6`

	_, err := m.DB.ExecContext(ctx, query,
		reservation.FirstName,
		reservation.LastName,
		reservation.Email,
		reservation.Phone,
		reservation.UpdatedAt,
		reservation.ID)

	if err != nil {
		return err
	}

	return nil
}

// DeleteReservation deletes a reservation from the database
func (m *postgresDBRepo) DeleteReservation(reservationID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations
			where id=$1`

	_, err := m.DB.ExecContext(ctx, query, reservationID)

	if err != nil {
		return err
	}

	return nil
}

// UpdateReservationAsProcessed updates a reservation in the database to be set as processed
func (m *postgresDBRepo) UpdateReservationAsProcessed(reservationID, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations
			set processed=$1
			where id=$2`

	_, err := m.DB.ExecContext(ctx, query, processed, reservationID)

	if err != nil {
		return err
	}

	return nil
}
