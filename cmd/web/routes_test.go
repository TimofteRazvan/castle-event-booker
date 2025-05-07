package main

import (
	"fmt"
	"testing"

	"github.com/TimofteRazvan/castle-event-booker/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		// correct, pass
	default:
		// error, should be pointer *chi.Mux
		t.Error(fmt.Sprintf("routes() return type %T is not *chi.Mux", v))
	}
}
