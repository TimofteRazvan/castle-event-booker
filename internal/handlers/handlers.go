package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TimofteRazvan/castle-event-booker/internal/config"
	"github.com/TimofteRazvan/castle-event-booker/internal/models"
	"github.com/TimofteRazvan/castle-event-booker/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository pattern type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Booking is the booking page handler
func (m *Repository) Booking(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "booking.page.tmpl", &models.TemplateData{})
}

// PostBooking posts the Booking page data
func (m *Repository) PostBooking(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// BookingJSON handles same-page request for booking and sends JSON response
func (m *Repository) BookingJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "available",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Knights is the knights page handler
func (m *Repository) Knights(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "knights.page.tmpl", &models.TemplateData{})
}

// Throne is the throne page handler
func (m *Repository) Throne(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "throne.page.tmpl", &models.TemplateData{})
}

// Banquet is the banquet page handler
func (m *Repository) Banquet(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "banquet.page.tmpl", &models.TemplateData{})
}

// Contact is the contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// MakeReservation is the make-reservation page handler
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}