package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TimofteRazvan/castle-event-booker/helpers"
	"github.com/TimofteRazvan/castle-event-booker/internal/config"
	"github.com/TimofteRazvan/castle-event-booker/internal/driver"
	"github.com/TimofteRazvan/castle-event-booker/internal/handlers"
	"github.com/TimofteRazvan/castle-event-booker/internal/models"
	"github.com/TimofteRazvan/castle-event-booker/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main app function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)
	fmt.Println("Starting mail listener on port :1025")
	listenForMail()

	// msg := models.MailData{
	// 	To:      "razvanhdt13@gmail.com",
	// 	From:    "corvinul_bookings@gmail.com",
	// 	Subject: "Your reservation was made!",
	// 	Content: "",
	// }

	// app.MailChan <- msg

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {

	// for storing non-primitives in session
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.Reservation{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when we're in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=postgres")
	if err != nil {
		log.Fatal("Could not connect to database")
	}
	log.Println("Connected to database")

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
