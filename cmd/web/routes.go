package main

import (
	"net/http"

	"github.com/TimofteRazvan/castle-event-booker/internal/config"
	"github.com/TimofteRazvan/castle-event-booker/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/knights", handlers.Repo.Knights)
	mux.Get("/throne", handlers.Repo.Throne)
	mux.Get("/banquet", handlers.Repo.Banquet)
	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/booking", handlers.Repo.Booking)
	mux.Post("/booking", handlers.Repo.PostBooking)
	mux.Post("/booking-json", handlers.Repo.BookingJSON)
	mux.Get("/book-room", handlers.Repo.BookRoom)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Post("/make-reservation", handlers.Repo.PostMakeReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login", handlers.Repo.PostLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		//mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/reservations-new", handlers.Repo.AdminReservationsNew)
		mux.Get("/reservations-all", handlers.Repo.AdminReservationsAll)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
		mux.Get("/reservations/{src}/{id}", handlers.Repo.AdminShowReservation)
	})

	return mux
}
