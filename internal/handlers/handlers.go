package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TimofteRazvan/castle-event-booker/internal/config"
	"github.com/TimofteRazvan/castle-event-booker/internal/driver"
	"github.com/TimofteRazvan/castle-event-booker/internal/forms"
	"github.com/TimofteRazvan/castle-event-booker/internal/models"
	"github.com/TimofteRazvan/castle-event-booker/internal/render"
	"github.com/TimofteRazvan/castle-event-booker/internal/repository"
	"github.com/TimofteRazvan/castle-event-booker/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository pattern type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Booking is the booking page handler
func (m *Repository) Booking(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "booking.page.tmpl", &models.TemplateData{})
}

// PostBooking posts the Booking page data
func (m *Repository) PostBooking(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// 01/02 03:04:05PM '06 -0700 (go date format)
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabilityByDateAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot search rooms by date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "no available rooms")
		http.Redirect(w, r, "/booking", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// ChooseRoom displays available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(splitURL[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// BookingJSON handles same-page request for booking and sends JSON response
func (m *Repository) BookingJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error parsing start date",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error parsing end date",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error converting string to int",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	available, err := m.DB.SearchAvailabilityByDateByRoomID(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		RoomID:    strconv.Itoa(roomID),
		StartDate: sd,
		EndDate:   ed,
	}

	out, _ := json.MarshalIndent(resp, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// BookRoom takes URL parameters and builds a reservation session variables,
// takes users to the make reservation page
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse id within URL")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "database query fail: cannot get room by id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var reservation models.Reservation
	reservation.RoomID = roomID
	reservation.StartDate = startDate
	reservation.EndDate = endDate
	reservation.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// Knights is the knights page handler
func (m *Repository) Knights(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "knights.page.tmpl", &models.TemplateData{})
}

// Throne is the throne page handler
func (m *Repository) Throne(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "throne.page.tmpl", &models.TemplateData{})
}

// Banquet is the banquet page handler
func (m *Repository) Banquet(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "banquet.page.tmpl", &models.TemplateData{})
}

// Contact is the contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// MakeReservation is the make-reservation page handler
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot find room based on the given id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", reservation)

	startDate := reservation.StartDate.Format("2006-01-02")
	endDate := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostMakeReservation is the POST make-reservation page / form handler
func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "form is invalid", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot insert reservation into the database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomRestriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(roomRestriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot insert room restriction into the database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// send notification to guest
	htmlMsg := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Reservation Confirmation</title>
		<style>
			body { font-family: Arial, sans-serif; background-color: #f4f4f4; color: #333; padding: 20px; }
			.container { max-width: 600px; margin: auto; background: #fff; padding: 20px; border-radius: 8px; }
			h2 { color: #2c3e50; }
			p { line-height: 1.6; }
			.footer { margin-top: 30px; font-size: 0.9em; color: #888; }
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Your Reservation was Confirmed!</h2>
			<p>Dear %s,</p>
			<p>Thank you for choosing Corvinul for your stay. We're happy to let you know that your reservation has been successfully confirmed.</p>
			<p><strong>Reservation Details:</strong></p>
			<ul>
				<li><strong>Check-in:</strong> %s</li>
				<li><strong>Check-out:</strong> %s</li>
			</ul>
			<p>If you have any questions or need to make changes to your reservation, feel free to contact us at <a href="mailto:corvinul_bookings@gmail.com">corvinul_bookings@gmail.com</a>.</p>
			<p>We look forward to welcoming you!</p>
			<p>Best regards,<br>The Corvinul Team</p>
			<div class="footer">
				<p>Corvinul Booking Service<br>corvinul_bookings@gmail.com</p>
			</div>
		</div>
	</body>
	</html>
	`, reservation.FirstName, reservation.StartDate.Format("January 2, 2006"), reservation.EndDate.Format("January 2, 2006"))

	msg := models.MailData{
		To:      reservation.Email,
		From:    "corvinul_bookings@gmail.com",
		Subject: "Your reservation was confirmed!",
		Content: htmlMsg,
	}

	m.App.MailChan <- msg

	htmlMsg = fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>New Reservation Notification</title>
		<style>
			body { font-family: Arial, sans-serif; background-color: #f4f4f4; color: #333; padding: 20px; }
			.container { max-width: 600px; margin: auto; background: #fff; padding: 20px; border-radius: 8px; }
			h2 { color: #2c3e50; }
			p { line-height: 1.6; }
			ul { padding-left: 20px; }
			.footer { margin-top: 30px; font-size: 0.9em; color: #888; }
		</style>
	</head>
	<body>
		<div class="container">
			<h2>New Reservation Received</h2>
			<p>Hello Razvan,</p>
			<p>You have received a new reservation for your property.</p>

			<p><strong>Reservation Details:</strong></p>
			<ul>
				<li><strong>Guest Name:</strong> %s %s</li>
				<li><strong>Email:</strong> %s</li>
				<li><strong>Check-in:</strong> %s</li>
				<li><strong>Check-out:</strong> %s</li>
			</ul>

			<p>Please make any necessary preparations and feel free to contact the guest if needed.</p>

			<p>Best regards,<br>The Corvinul Booking Team</p>

			<div class="footer">
				<p>This is an automated message sent from Corvinul Booking Service.</p>
			</div>
		</div>
	</body>
	</html>
`, reservation.FirstName, reservation.LastName, reservation.Email, reservation.StartDate.Format("January 2, 2006"), reservation.EndDate.Format("January 2, 2006"))

	msg = models.MailData{
		To:      "corvinul_bookings@gmail.com",
		From:    "corvinul_bookings@gmail.com",
		Subject: "A reservation has been made!",
		Content: htmlMsg,
	}

	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary displays the details of the posted reservation
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	startDate := reservation.StartDate.Format("2006-01-02")
	endDate := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}
