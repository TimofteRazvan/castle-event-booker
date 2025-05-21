package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/TimofteRazvan/castle-event-booker/internal/config"
)

var app *config.AppConfig

// Sets up the app config for helpers
func NewHelpers(appConfig *config.AppConfig) {
	app = appConfig
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error code ", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}
